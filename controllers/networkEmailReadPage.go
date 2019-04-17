package controllers

import (
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../models"
    "../crypto"
    "../settings"
)

func networkEmailReadPage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    var data = struct {
        Auth bool
        Login string
        UserHash string
        Emails []models.Email
    } {}
    
    var emails []models.Email
    var email models.Email

    if r.Method == "POST" {
        r.ParseForm()
        if _, ok := r.Form["delete"]; ok {
            _, err := settings.DataBase.Exec("DELETE FROM Email WHERE id=$1", r.FormValue("id"))
            utils.CheckError(err)
        }
    }

    if r.URL.Path == "/network/email/read/" {
        rows, err := settings.DataBase.Query("SELECT Id, Title, User, Date FROM Email")
        utils.CheckError(err)

        for rows.Next() {
            err = rows.Scan(
                &email.Id,
                &email.Title,
                &email.User,
                &email.Date,
            )
            utils.CheckError(err)
            crypto.DecryptEmail(settings.User.Password, &email)
            emails = append(emails, email)
        }

        rows.Close()

        data.Auth = true
        data.Login = settings.User.Login
        data.Emails = emails

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_email_read.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var slice = strings.Split(strings.TrimPrefix(r.URL.Path, "/network/email/read/"), "/")
    var username = slice[0]

    switch len(slice) {
        case 1: 
            rows, err := settings.DataBase.Query(
                "SELECT Id, Title, User, Date FROM Email WHERE User=$1", 
                username,
            )
            utils.CheckError(err)

            for rows.Next() {
                err = rows.Scan(
                    &email.Id,
                    &email.Title,
                    &email.User,
                    &email.Date,
                )
                utils.CheckError(err)
                crypto.DecryptEmail(settings.User.Password, &email)
                emails = append(emails, email)
            }

            rows.Close()

            data.Auth = true
            data.Login = settings.User.Login
            data.UserHash = username
            data.Emails = emails

            tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_email_read_X.html")
            utils.CheckError(err)
            tmpl.Execute(w, data)
            return

        case 2: 
            var id = slice[1]
            rows, err := settings.DataBase.Query(
                "SELECT * FROM Email WHERE User=$1 AND Id=$2", 
                username, 
                id,
            )
            utils.CheckError(err)

            for rows.Next() {
                err = rows.Scan(
                    &email.Id,
                    &email.Title,
                    &email.Body,
                    &email.User,
                    &email.Date,
                )
                utils.CheckError(err)
                crypto.DecryptEmail(settings.User.Password, &email)
                emails = append(emails, email)
            }

            rows.Close()

            data.Auth = true
            data.Login = settings.User.Login
            data.UserHash = username
            data.Emails = emails

            tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_email_read_X_Y.html")
            utils.CheckError(err)
            tmpl.Execute(w, data)
            return

        default: 
            redirectTo("page404", w, r)
            return
    }
}
