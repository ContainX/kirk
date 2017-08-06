package routes

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ContainX/kirk/config"
	"github.com/ContainX/kirk/slackTeam"
	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
)

func Auth(router *gin.Engine) {
	// Some statics for handling oauth
	authorizedRoute := "/authorized"
	authorizedRouteRedirectUri := "http://" + os.Getenv("PUBLIC_HOST") + authorizedRoute
	oauthScopes := []string{
		"channels:read",
		"chat:write:bot", // or bot?
		"team:read",
		"channels:history",
		"groups:history",
		"im:history",
		"im:read",
		"im:write",
		"mpim:history",
		"bot",
	}

	queryParams := url.Values{}
	queryParams.Set("client_id", os.Getenv("SLACK_CLIENT_ID"))
	queryParams.Set("scope", strings.Join(oauthScopes, " "))
	queryParams.Set("redirect_uri", authorizedRouteRedirectUri)

	// 1. User visits initial authorize route and is redirected to Slack
	router.GET("/authorize", func(context *gin.Context) {
		context.Redirect(308, "https://slack.com/oauth/authorize?"+queryParams.Encode())
	})

	// 2. Once user authorizes app on Slack's end, Slack redirects user back to Kirk.
	router.GET(authorizedRoute, func(context *gin.Context) {
		authCode := context.Query("code")

		fmt.Println("Authorized received.  Code:", authCode)
		oAuthResponse, err := slack.GetOAuthResponse(
			os.Getenv("SLACK_CLIENT_ID"),
			os.Getenv("SLACK_CLIENT_SECRET"),
			authCode,
			authorizedRouteRedirectUri,
			false,
		)

		if err != nil {
			fmt.Println("Error getting oath token")
			context.HTML(500, "oauthError.html", struct{}{})
		} else {
			fmt.Println("Saving token")
			config.AddNewTeam(oAuthResponse.Bot.BotAccessToken)
			fmt.Println("Starting new bot instance")
			slackTeam.Watch(oAuthResponse.Bot.BotAccessToken)
			context.HTML(200, "oauthSuccess.html", struct{}{})
		}
	})
}
