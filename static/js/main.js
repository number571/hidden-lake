import { RoutesData } from "./routes.js";

const apihost = `${http}${host}/api/`;
const routes = {
    home: 0,
    about: 1,
    login: 2,
    signup: 3,
    account: 4,
    friends: 5,
    archive: 6,
    archivefile: 7,
    network: 8,
    settings: 9,
    clients: 10,
    globalchatlist: 11,
    globalchat: 12,
    chatnull: 13,
    chat: 14,
    emailnull: 15,
    email: 16,
    client: 17,
    clientarchive: 18,
    clientarchivefile: 19,
    hdvwhuhjj: 20,
    notfound: 21,
};

const f = async(url, method = "GET", data = null, token = null) => {
    method = method.toLocaleUpperCase()
    let fullurl = `${apihost}${url}`;
    let options = {
        method: method, 
        headers: {
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`,
        },
    };

    switch(method) {
        case "PUT":
            delete options.headers["Content-Type"];
            options.body = data;
            break;
        case "POST": case "PATCH": case "DELETE":
            options.body = JSON.stringify(data);
            break;
    }

    const res = await fetch(fullurl, options);
    return await res.json();
};

const router = new VueRouter({
    mode: "history",
    routes: RoutesData,
});

router.beforeEach((to, from, next) => {
    document.title = to.meta.title;
    next();
});

Vue.component("notfound", {
    template: "",
    created: function() {
        // Redirect outside the app using plain old javascript
        window.location.href = "/notfound";
    },
});

const app = new Vue({
    delimiters: ['${', '}'],
    el: "main",
    router: router,
    data: {
        userdata: {
            username: null,
            hashname: null,
            password: null,
            statef2f: null,
            password_repeat: null,
            private_key: null,
            connects: [],
        },
        authdata: {
            token: null,
            username: null,
        },
        conndata: {
            in_chat: null,
            connected: null,
            hidden: null,
            throwclient: null,
            address: null,
            hashname: null,
            public_key: null,
            certificate: null,
            message: null,
        },
        filelist: [],
        filedata: {
            name: null,
            hash: null,
            path: null,
            size: null,
            encr: false,
        },
        emaillist: [],
        emaildata: {
            info: {
                incoming: null,
                time: null,
            },
            email: {
                head: {
                    sender: {
                        public_key: null,
                        hashname: null,
                    },
                    receiver: null,
                    session: null,
                },
                body: {
                    data: {
                        title: null,
                        message: null,
                    },
                    desc: {
                        rand: null,
                        hash: null,
                        sign: null,
                        nonce: null,
                        difficulty: null,
                    },
                },
            },
        },
        netdata: {
            message: null,
            list: [],
            chat: {
                companion: null,
                messages: [],
            },
        },
        socket: null,
        switcher: null,
        message: {
            wait: null,
            curr: null,
            desc: null,
        },
        checked: false,
        opened: null,
    },
    methods: {
        async login(username, password) {
            let obj = {
                username: username, 
                password: password,
            };
            let res = await f("login", "POST", obj);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            localStorage.setItem("token", res.token);
            localStorage.setItem("username", username);

            this.authdata.token = localStorage.getItem("token");
            this.authdata.username = localStorage.getItem("username");

            this.nulldata();

            this.message.wait = "Login success";
            this.message.desc = "success";

            this.opened = RoutesData[routes.home].name;
            this.$router.push(RoutesData[routes.home]);
        },
        async signup(username, password, password_repeat, private_key) {
            if (password !== password_repeat) {
                this.message.curr = "Passwords not equal";
                this.message.desc = "danger";
                return;
            }
            if (
                username.length < 6 || password.length < 8 ||
                username.length > 128 || password.length > 1024
            ) {
                this.message.curr = "Username needs [6-128] ch and password needs [8-1024] ch";
                this.message.desc = "danger";
                return;
            }
            let obj = {
                username: username,
                password: password,
                password_repeat: password_repeat,
                private_key: private_key,
            };
            let res = await f("signup", "POST", obj);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.nulldata();

            this.message.wait = "Signup success";
            this.message.desc = "success";

            this.opened = RoutesData[routes.login].name;
            this.$router.push(RoutesData[routes.login]);
        },
        async logout() {
            let res = await f("logout", "POST", null, this.authdata.token);
            this.nulldata();
            this.nullauth();
        },
        async account() {
            let res = await f("account", "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.conndata.hashname = res.hashname;
            this.conndata.address = res.address;
            this.conndata.public_key = res.public_key;
            this.conndata.certificate = res.certificate;
        },
        async viewkey(username, password) {
            let obj = {
                username: username,
                password: password,
            };
            let res = await f("account", "POST", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.userdata.private_key = res.private_key;
        },
        async deluser(username, password) {
            let obj = {
                username: username,
                password: password,
            };
            let res = await f("account", "DELETE", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.nullauth();
            this.message.wait = "Delete success";
            this.message.desc = "success";
            this.opened = RoutesData[routes.login].name;
            this.$router.push(RoutesData[routes.login]);
        },
        async allconnects() {
            let res = await f("account/connects", "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.userdata.connects = res.connects;
        },
        async currconnects() {
            let res = await f("account/connects", "PATCH", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.userdata.connects = res.connects;
        },
        async email(hash) {
            let res = await f(`network/email/${hash}`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            switch(hash) {
                case "": case "null": case "undefined":
                    this.emaillist = res.emails;
                    break;
                default:
                    this.emaildata = res.email;
                    break;
            }
        },
        async emailupdate() {
            let res = await f(`network/email/`, "PATCH", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Update success";
            this.message.desc = "success";
            this.email("null");
        },
        async emailsend(public_key, title, message) {
            let obj = {
                public_key: public_key,
                title: title,
                message: message,
            };
            let res = await f(`network/email/`, "POST", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.email("null");
            this.message.curr = "Send success";
            this.message.desc = "success";
        },
        async emaildel(hash) {
            let res = await f(`network/email/`, "DELETE", {emailhash: hash}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.email("null");
            if (this.opened === RoutesData[routes.emailnull].name) {
                this.message.curr = "Delete success";
                this.message.desc = "success";
            } else {
                this.message.wait = "Delete success";
                this.message.desc = "success";
                this.opened = RoutesData[routes.emailnull].name;
                this.$router.push(RoutesData[routes.emailnull]);
            }
        },
        async chat(name) {
            let res = await f(`network/chat/${name}`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.netdata.list = res.netdata.list;
            this.netdata.chat = res.netdata.chat;
            if (this.socket != null) {
                this.socket.close(1000, "new socket");
            }
            this.socket = new WebSocket(`${ws}${host}/ws/network`);
            this.socket.onopen = () => {
                this.socket.send(JSON.stringify({
                    token: this.authdata.token, 
                    option: "private",
                }));
            }
            this.socket.onmessage = (e) => {
                let obj = JSON.parse(e.data);
                for (let i = 0; i < this.netdata.list.length; ++i) {
                    if (this.netdata.list[i].companion == name) {
                        this.netdata.list[i].message.name = obj.comp.from;
                        this.netdata.list[i].message.text = obj.text;
                        this.netdata.list[i].message.time = obj.time;
                        break;
                    }
                }
                if (obj.comp.from != name && obj.comp.to != name) {
                    return;
                }
                this.netdata.chat.messages.push({
                    name: obj.comp.from,
                    text: obj.text,
                    time: obj.time,
                });
            }
            this.socket.onerror = (e) => {
                // console.debug(e);
            }
            this.socket.onclose = (e) => {
                // console.debug("closed");
            }
        },
        async sendmsg(companion, message) {
            let obj = {
                hashname: companion,
                message: message,
            };
            let res = await f("network/chat/", "POST", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.netdata.message = null;
            this.message.curr = null;
        },
        async delchat(hashname, username, password) {
            let obj = {
                hashname: hashname,
                username: username,
                password: password,
            };
            let res = await f("network/chat/", "DELETE", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }

            this.netdata.list.splice(this.netdata.list.indexOf(this.conndata.hashname), 1);
            this.netdata.chat.companion = null;
            this.netdata.chat.messages = [];

            this.message.wait = "Delete chat success";
            this.message.desc = "success";
            this.opened = RoutesData[routes.network].name;
            this.$router.push(RoutesData[routes.network]);
        },
        async globalchat(hashname) {
            let res = await f(`network/chat/global/${hashname}`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.netdata.list = res.list;
            switch (hashname) {
                case "": case "null": case "undefined":
                    return
            }
            this.netdata.chat = res.chat;
            if (this.socket != null) {
                this.socket.close(1000, "new socket");
            }
            this.socket = new WebSocket(`${ws}${host}/ws/network`);
            this.socket.onopen = () => {
                this.socket.send(JSON.stringify({
                    token: this.authdata.token, 
                    option: "group",
                }));
            }
            this.socket.onmessage = (e) => {
                let obj = JSON.parse(e.data);
                if (obj.comp.to != hashname) {
                    return;
                }
                this.netdata.chat.messages.push({
                    name: obj.comp.from,
                    text: obj.text,
                    time: obj.time,
                });
            }
            this.socket.onerror = (e) => {
                // console.debug(e);
            }
            this.socket.onclose = (e) => {
                // console.debug("closed");
            }
        },
        async exitglobalchat(hashname) {
            let obj = {
                hashname: hashname,
                option: 'exit',
            };
            let res = await f(`network/chat/global/${hashname}`, "PATCH", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Exit from chat success";
            this.message.desc = "success";
        },
        async joinglobalchat(hashname) {
            let obj = {
                hashname: hashname,
                option: 'join',
            };
            let res = await f(`network/chat/global/${hashname}`, "PATCH", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Join to chat success";
            this.message.desc = "success";
        },
        async sendglobalchat(founder, message) {
            let obj = {
                hashname: founder,
                message: message,
            };
            let res = await f(`network/chat/global/`, "POST", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.netdata.message = null;
            this.message.curr = null;
        },
        async delglobalchat(hashname, username, password, option) {
            let obj = {
                hashname: hashname,
                username: username,
                password: password,
                password_repeat: option,
            };
            let res = await f("network/chat/global/", "DELETE", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.wait = "Delete chat success";
            this.message.desc = "success";

            this.globalchat("null");
            this.opened = RoutesData[routes.globalchatlist].name;
            this.$router.push(RoutesData[routes.globalchatlist]);
        },
        async delclient(hashname, username, password) {
            let obj = {
                hashname: hashname,
                username: username,
                password: password,
            };
            let res = await f("account/connects", "DELETE", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }

            this.message.wait = "Delete client success";
            this.message.desc = "success";
            this.opened = RoutesData[routes.network].name;
            this.$router.push(RoutesData[routes.network]);
        },
        async client(hashname) {
            let res = await f(`network/client/${hashname}`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.conndata.in_chat = res.info.in_chat;
            this.conndata.connected = res.info.connected;
            this.conndata.hidden = res.info.hidden;
            this.conndata.throwclient = res.info.throwclient;
            this.conndata.address = res.info.address;
            this.conndata.hashname = res.info.hashname;
            this.conndata.public_key = res.info.public_key;
            this.conndata.certificate = res.info.certificate;
        },
        async connect(address, certificate, public_key) {
            this.message.curr = "Please wait a few seconds";
            this.message.desc = "warning";
            let obj = {
                address: address, 
                certificate: certificate,
                public_key: public_key,
            };
            let res = await f("network/client/", "POST", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Connection success"
            this.message.desc = "success";
        },
        async disconnect(hashname) {
            this.message.curr = "Please wait a few seconds";
            this.message.desc = "warning";
            let res = await f("network/client/", "DELETE", {hashname: hashname}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Disconnection success"
            this.message.desc = "success";
        },
        async archivelist(hashname) {
            this.nullfile();
            if (hashname == '') {
                let res = await f(`account/archive/`, "GET", null, this.authdata.token);
                if (res.state) {
                    this.message.curr = res.state;
                    this.message.desc = "warning";
                    return;
                }
                this.filelist = res.files;
                return;
            }
            let res = await f(`network/client/${hashname}/archive/`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.conndata.hashname = hashname;
            this.filelist = res.files;
        },
        async archivefile(hashname, filehash) {
            this.nullfile();
            if (hashname == '') {
                let res = await f(`account/archive/${filehash}`, "GET", null, this.authdata.token);
                if (res.state) {
                    this.message.curr = res.state;
                    this.message.desc = "warning";
                    return;
                }
                this.filedata.name = res.files[0].name;
                this.filedata.hash = res.files[0].hash;
                this.filedata.path = res.files[0].path;
                this.filedata.size = res.files[0].size;
                this.filedata.encr = res.files[0].encr;
                return;
            }
            let res = await f(`network/client/${hashname}/archive/${filehash}`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.conndata.hashname = hashname;
            this.filedata.name = res.files[0].name;
            this.filedata.hash = res.files[0].hash;
            this.filedata.path = res.files[0].path;
            this.filedata.size = res.files[0].size;
            this.filedata.encr = res.files[0].encr;
        },
        async installfile(hashname, filehash) {
            this.message.curr = "Please wait a few seconds";
            this.message.desc = "warning";
            let res = await f(`network/client/${hashname}/archive/`, "POST", {filehash: filehash}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Download success";
            this.message.desc = "success";
        },
        async downloadfile(filehash) {
            var win = window.open(`${http}${host}/static/archive/${filehash}?token=${encodeURIComponent(this.authdata.token)}`, '_blank');
            win.focus();
            return;
        },
        async deletefile(filehash) {
            let obj = {
                hashname: filehash,
                username: username,
                password: password,
            };
            let res = await f(`account/archive/`, "DELETE", obj, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.wait = "Delete success";
            this.message.desc = "success";

            this.archivelist('');
            this.opened = RoutesData[routes.archive].name;
            this.$router.push(RoutesData[routes.archive]);
        },
        async uploadfile(checked) {
            this.message.curr = "Please wait a few seconds";
            this.message.desc = "warning";
            const formData = new FormData();
            const fileField = document.querySelector('#uploadfile');
            formData.append("encryptmode", checked);
            formData.append("uploadfile", fileField.files[0]);
            let res = await f(`account/archive/`, "PUT", formData, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = `Upload success; Hash: ${res.filehash}`;
            this.message.desc = "success";
            this.archivelist('');
            this.opened = RoutesData[routes.archive].name;
            this.$router.push(RoutesData[routes.archive]);
        },
        async findconnect(public_key) {
            let res = await f(`network/client/`, "PATCH", {public_key: public_key}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.currconnects();
            this.message.curr = "Connection success";
            this.message.desc = "success";
        },
        async getfriends() {
            let res = await f(`account/friends`, "GET", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "warning";
                return;
            }
            this.userdata.statef2f = res.statef2f;
            this.userdata.connects = res.friends;
        },
        async delfriend(hashname) {
            let res = await f(`account/friends`, "DELETE", {hashname: hashname}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Delete friend success";
            this.message.desc = "success";
            this.getfriends();
        },
        async addfriend(hashname) {
            let res = await f(`account/friends`, "PATCH", {hashname: hashname}, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.message.curr = "Append friend success";
            this.message.desc = "success";
            this.getfriends();
        },
        async turnf2f() {
            let res = await f(`account/friends`, "POST", null, this.authdata.token);
            if (res.state) {
                this.message.curr = res.state;
                this.message.desc = "danger";
                return;
            }
            this.userdata.statef2f = res.statef2f;
            this.message.curr = `Turn ${res.statef2f ? "ON" : "OFF"} success`;
            this.message.desc = "success";
        },
        selectText(element) {
            var range;
            if (document.selection) {
                range = document.body.createTextRange();
                range.moveToElementText(element);
                range.select();
            } else if (window.getSelection) {
                range = document.createRange();
                range.selectNode(element);
                window.getSelection().removeAllRanges();
                window.getSelection().addRange(range);
            }
        },
        savepublic() {
            this.selectText(this.$refs.publickey);
            let res = document.execCommand("copy");
            if (!res) {
                this.message.curr = "Public key not copied to clipboard"
                this.message.desc = "danger";
                return
            }
            this.message.curr = "Public key copied to clipboard successfully"
            this.message.desc = "success";
        },
        savecertificate() {
            this.selectText(this.$refs.certificate);
            let res = document.execCommand("copy");
            if (!res) {
                this.message.curr = "Public key not copied to clipboard"
                this.message.desc = "danger";
                return
            }
            this.message.curr = "Certificate copied to clipboard successfully"
            this.message.desc = "success";
        },
        nullauth() {
            this.authdata.token = null;
            this.authdata.username = null;
            localStorage.removeItem("token");
            localStorage.removeItem("username");
        },
        nullemail() {
            this.emaillist = [];
            this.emaildata = {
                info: {
                    incoming: null,
                    time: null,
                },
                email: {
                    head: {
                        sender: {
                            public_key: null,
                            hashname: null,
                        },
                        receiver: null,
                        session: null,
                    },
                    body: {
                        data: {
                            title: null,
                            message: null,
                        },
                        desc: {
                            rand: null,
                            hash: null,
                            sign: null,
                            nonce: null,
                            difficulty: null,
                        },
                    },
                },
            };
        },
        nullfile() {
            this.filedata.name = null;
            this.filedata.hash = null;
            this.filedata.path = null;
            this.filedata.size = null;
            this.filedata.encr = false;
            this.filelist = [];
        },
        nullconn() {
            this.conndata.in_chat = null;
            this.conndata.connected = null;
            this.conndata.hidden = null;
            this.conndata.throwclient = null;
            this.conndata.hashname = null;
            this.conndata.address = null;
            this.conndata.public_key = null;
            this.conndata.certificate = null;
            this.conndata.message = null;
        },
        nulldata() {
            this.userdata.username = null;
            this.userdata.password = null;
            this.userdata.password_repeat = null;
            this.userdata.private_key = null;
            this.userdata.connects = [];
            this.userdata.statef2f = null;
        },
        nullcurr(page) {
            this.switcher = null;
            this.checked = false;
            this.message.curr = null;
            switch (page) {
                case RoutesData[routes.login].name:
                case RoutesData[routes.signup].name:
                    this.nulldata();
                    break;
            }
        },
        setswitch(name) {
            this.switcher = (this.switcher === name) ? null : name;
        },
    },
    mounted() {
        let token = localStorage.getItem("token");
        if (token) {
            this.authdata.token = localStorage.getItem("token");
            this.authdata.username = localStorage.getItem("username");
        }
        switch (this.$route.name) {
            case RoutesData[routes.settings].name: this.currconnects(); break;
            case RoutesData[routes.friends].name: this.getfriends(); break;
            case RoutesData[routes.account].name: this.account(); break;
            case RoutesData[routes.archive].name: this.archivelist(''); break;
            case RoutesData[routes.archivefile].name: this.archivefile('', this.$route.params.id); break;
            case RoutesData[routes.chat].name: this.chat(this.$route.params.id); break;
            case RoutesData[routes.chatnull].name: this.chat("null"); break;
            case RoutesData[routes.globalchat].name: this.globalchat(this.$route.params.id); break;
            case RoutesData[routes.globalchatlist].name: this.globalchat("null"); break;
            case RoutesData[routes.email].name: this.email(this.$route.params.id); break;
            case RoutesData[routes.emailnull].name: this.email("null"); break;
            case RoutesData[routes.client].name: this.client(this.$route.params.id); break;
            case RoutesData[routes.clients].name: this.allconnects(); break;
            case RoutesData[routes.clientarchive].name: this.archivelist(this.$route.params.id); break;
            case RoutesData[routes.clientarchivefile].name: this.archivefile(this.$route.params.id0, this.$route.params.id1); break;
        }
        this.opened = this.$route.name;
    },
    updated() {
        switch (this.opened) {
            case RoutesData[routes.chat].name: 
            case RoutesData[routes.globalchat].name:
                this.$nextTick(() => {
                    var bottomChat = this.$refs.bottomChat;
                    bottomChat.scrollTop = bottomChat.scrollHeight;
                });
                break;
        }
    },
    watch: {
        '$route' (to, from) {
            this.nullcurr(to.name);
            this.opened = to.name;
            if (this.message.wait != null) {
                this.message.curr = this.message.wait;
                this.message.wait = null;
            }
        },
    },
});
