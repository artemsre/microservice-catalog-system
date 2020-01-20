package main

import (
	"database/sql"
	"fmt"
	oidc "github.com/coreos/go-oidc"
	"github.com/foolin/echo-template"
	"github.com/go-session/echo-session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

type Service struct {
	Id           string
	Name         string
	Type         string
	Level        string
	Status       string
	Team         string
	Product      string
	DocUrl       string
	SvcUrl       string
	DashboardUrl string
	Newrelic_id  string
	Sentry_id    string
}

func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Fatal(err)
	}
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://127.0.0.1:9090/auth/google/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	servcatid := "AreYouAHackerOrTheHacker"
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.Redirect(200, "/")
	}
	// Echo instance
	e := echo.New()
	e.Static("/static", "static")
	e.Use(echosession.New())
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = echotemplate.Default()

	// Routes
	e.GET("/", func(c echo.Context) error {
		session := echosession.FromContext(c)
		sessionId, ok := session.Get("servcatid")
		if !ok {
			return c.Redirect(http.StatusFound, config.AuthCodeURL(servcatid))
		}
		user, ok := session.Get("username")
		if !ok {
			user = "Null"
		}
		if sessionId != servcatid {
			return c.Redirect(http.StatusFound, config.AuthCodeURL(servcatid))
		}
		services, err := listServices()
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.Render(http.StatusOK, "dyn.html",
			echo.Map{"user": user, "total_services": 120, "internal_services": 90,
				"third_services": 30, "databases": 0, "services": services})
	})

	e.GET("/auth/google/callback", func(c echo.Context) error {
		if c.QueryParam("state") != servcatid {
			return echo.NewHTTPError(http.StatusBadRequest, "state did not match")
		}

		oauth2Token, err := config.Exchange(ctx, c.QueryParam("code"))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
		}

		userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get userinfo: "+err.Error())
		}

		name := userInfo.Email
		print("name=" + name)
		if userInfo.EmailVerified == false {
			return echo.NewHTTPError(403, "Email: "+name+" is not verufed at google")
		}
		valid_domain := os.Getenv("VALID_DOMAIN")
		if strings.HasSuffix(name, valid_domain) {
			session := echosession.FromContext(c)
			session.Set("servcatid", servcatid)
			session.Set("username", name)
			session.Save()
			return c.Redirect(302, "/")
		}
		return echo.NewHTTPError(403, "Forbidden for user: "+name)
	})

	e.GET("/edit", func(c echo.Context) error {
		if checkSession(c) == false {
			c.Redirect(302, "/")
		}
		session := echosession.FromContext(c)
		user, ok := session.Get("username")
		if !ok {
			user = "Null"
		}
		svid := c.QueryParam("id")
		dberr := c.QueryParam("err")
		service, err := getService(svid)
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		teams, err := listTeams()
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		products, err := listProducts()
		if err != nil {
			return c.String(http.StatusOK, err.Error())
		}
		return c.Render(http.StatusOK, "edit.html",
			echo.Map{"user": user, "service": service, "products": products, "teams": teams, "error": dberr})
	})

	e.POST("/edit", func(c echo.Context) error {
		if checkSession(c) == false {
			c.Redirect(302, "/")
		}
		var s Service
		s.Id = c.FormValue("id")
		s.Name = c.FormValue("name")
		s.Type = c.FormValue("type")
		s.Level = c.FormValue("level")
		s.Status = c.FormValue("status")
		s.Team = c.FormValue("team")
		s.Product = c.FormValue("product")
		s.DocUrl = c.FormValue("docurl")
		s.SvcUrl = c.FormValue("svcurl")
		s.DashboardUrl = c.FormValue("dashboardurl")
		s.Newrelic_id = c.FormValue("newrelic_id")
		s.Sentry_id = c.FormValue("sentry_id")
		fmt.Printf("%+v\n", s)
		err := setService(s)
		if err != nil {
			return c.Redirect(302, "/edit?id="+c.FormValue("id")+"&err="+err.Error())
		}
		return c.Redirect(302, "/?msg=Success")

	})

	// Start server
	e.Logger.Fatal(e.Start("127.0.0.1:9090"))
}

func listServices() ([]Service, error) {
	services := []Service{}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return services, err
	}
	defer db.Close()
	rows, err := db.Query("select s.id,s.name,COALESCE(s.service_type,''),COALESCE(s.service_level,''),COALESCE(s.service_status,''),COALESCE(p.name,''),COALESCE(t.name,'') from services s left join teams t on t.team_id=s.team_id left join products p on p.product_id=s.product_id;")
	if err != nil {
		return services, err
	}
	for rows.Next() {
		var s Service
		err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.Level, &s.Status, &s.Product, &s.Team)
		if err != nil {
			return services, err
		}
		services = append(services, s)
	}
	return services, nil
}
func listProducts() ([]string, error) {
	products := []string{}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return products, err
	}
	defer db.Close()
	rows, err := db.Query("select name from products;")
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return products, err
		}
		products = append(products, s)
	}
	return products, nil
}
func listTeams() ([]string, error) {
	teams := []string{}
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return teams, err
	}
	defer db.Close()
	rows, err := db.Query("select name from teams;")
	if err != nil {
		return teams, err
	}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return teams, err
		}
		teams = append(teams, s)
	}
	return teams, nil
}

func getService(serviceid string) (Service, error) {
	var s Service
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return s, err
	}
	defer db.Close()
	rows, err := db.Query("select s.id,s.name,COALESCE(s.service_type,''),COALESCE(s.service_level,''),COALESCE(s.service_status,''),COALESCE(p.name,''),COALESCE(t.name,''), COALESCE(s.sentry_id,''),COALESCE(s.newrelic_id,''),COALESCE(s.dashboard_url,''), COALESCE(s.docurl,''), COALESCE(s.svcurl,'') from services s left join teams t on t.team_id=s.team_id left join products p on p.product_id=s.product_id where id=$1;", serviceid)
	if err != nil {
		return s, err
	}
	for rows.Next() {
		err := rows.Scan(&s.Id, &s.Name, &s.Type, &s.Level, &s.Status, &s.Product, &s.Team, &s.Sentry_id, &s.Newrelic_id, &s.DashboardUrl, &s.DocUrl, &s.SvcUrl)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func setService(service Service) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStatement := `UPDATE services set name=$1,service_type=$2, service_level=$3, service_status=$4, sentry_id=$5, newrelic_id=$6, dashboard_url=$7, product_id=(select product_id from products where name=$8 limit 1), team_id=(select team_id from teams where name=$9 limit 1), docurl=$11, svcurl=$12 where id=$10;`
	_, err = db.Exec(sqlStatement, service.Name, service.Type, service.Level, service.Status, service.Sentry_id, service.Newrelic_id, service.DashboardUrl, service.Product, service.Team, service.Id, service.DocUrl, service.SvcUrl)
	if err != nil {
		return err
	}
	return nil
}

func checkSession(c echo.Context) bool {
	session := echosession.FromContext(c)
	sessionId, ok := session.Get("servcatid")
	if !ok {
		return false
	}
	if sessionId != "AreYouAHackerOrTheHacker" {
		return false
	}
	return true
}
