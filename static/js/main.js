import { RoutesData } from "./routes.js";

/*
$ openssl genrsa -out https-server.key 2048
$ openssl ecparam -genkey -name secp384r1 -out https-server.key
$ openssl req -new -x509 -sha256 -key https-server.key -out https-server.crt -days 3650
*/

const apihost = `${http}${host}/api/`;
const routes = {
    home: 0,
    about: 1,
    login: 2,
    signup: 3,
    account: 4,
    network: 5,
    settings: 6,
    client: 7,
    notfound: 8,
};

const fapi = async(url, method = "GET", data = null, token = null) => {
    method = method.toLocaleUpperCase()
    let fullurl = `${apihost}${url}`;
    let options = {url, method, headers: {}};
    options.headers["Accept"] = "application/json";
    options.headers["Content-Type"] = "application/json";
    if (token) {
        options.headers["Authorization"] = `Bearer ${token}`;
    }
    if (["POST", "PUT", "PATCH"].includes(method)) {
        options.body = JSON.stringify(data);
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

const app = new Vue({
    delimiters: ['${', '}'],
    el: "main",
    router: router,
    data: {
        userdata: {
            username: null,
            password: null,
            password_repeat: null,
            private_key: null,
        },
        authdata: {
            token: null,
            username: null,
        },
        conndata: {
            connected: null,
            address: null,
            hashname: null,
            public_key: null,
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
        message: null,
        opened: null,
    },
    methods: {
        async login() {
            let res = await fapi("login", "POST", this.userdata);
            if (res.state) {
                this.message = res.state;
                return;
            }
            localStorage.setItem("token", res.token);
            localStorage.setItem("username", this.userdata.username);

            this.authdata.token = localStorage.getItem("token");
            this.authdata.username = localStorage.getItem("username");
            this.opened = RoutesData[routes.home].name;
            this.$router.push(RoutesData[routes.home]);
        },
        async signup() {
            if (this.userdata.password !== this.userdata.password_repeat) {
                this.message = "Passwords not equal";
                return;
            }
            if (
                this.userdata.username.length < 6 || this.userdata.password.length < 6 ||
                this.userdata.username.length > 64 || this.userdata.password.length > 128
            ) {
                this.message = "Username needs [6-64] ch and password needs [6-128] ch";
                return;
            }
            let res = await fapi("signup", "POST", this.userdata);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.nulldata();
            this.opened = RoutesData[routes.login].name;
            this.message = "Signup success";
        },
        async logout() {
            let res = await fapi("logout", "POST", null, this.authdata.token);
            this.nullauth();
        },
        async account() {
            let res = await fapi("account", "POST", null, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.conndata.hashname = res.hashname;
            this.conndata.address = res.address;
            this.conndata.public_key = res.public_key;
        },
        async viewkey() {
            let res = await fapi("viewkey", "POST", this.userdata, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.userdata.private_key = res.private_key;
        },
        async deluser() {
            let res = await fapi("deluser", "POST", this.userdata, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.nullauth();
            this.opened = RoutesData[routes.login].name;
        },
        async network(name) {
            let res = await fapi("network", "POST", {hashname: name}, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.netdata.list = res.netdata.list;
            this.netdata.chat = res.netdata.chat;
            if (this.socket != null) {
                this.socket.close(1000, "new socket");
            }
            this.socket = new WebSocket(`${ws}${host}/ws/network`);
            this.socket.onopen = () => {
                this.socket.send(JSON.stringify({token: this.authdata.token}));
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
        async client(name) {
            let res = await fapi("client", "POST", {hashname: name}, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.conndata.connected = res.connected;
            this.conndata.address = res.address;
            this.conndata.hashname = res.hashname;
            this.conndata.public_key = res.public_key;
        },
        async connect() {
            this.message = "Please wait a few seconds";
            let res = await fapi("connect", "POST", this.conndata, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.message = "Connection success"
        },
        async disconnect() {
            this.message = "Please wait a few seconds";
            let res = await fapi("disconnect", "POST", {hashname: this.conndata.hashname}, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.conndata.connected = false;
            this.message = "Disconnection success"
        },
        async sendmsg() {
            let obj = {
                hashname: this.netdata.chat.companion,
                message: this.netdata.message,
            };
            let res = await fapi("sendmsg", "POST", obj, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
        },
        async delchat() {
            let obj = {
                hashname: this.conndata.hashname,
                username: this.userdata.username,
                password: this.userdata.password,
            };
            let res = await fapi("delchat", "POST", obj, this.authdata.token);
            if (res.state) {
                this.message = res.state;
                return;
            }
            this.message = "Chat deleted success";
        },
        async keycheck(e) {
            if (e.keyCode == 13) { // Enter
                this.sendmsg();
                this.netdata.message = "";
            }
        },
        nullauth() {
            this.authdata = {
                token: null,
                username: null,
            };
            localStorage.removeItem("token");
            localStorage.removeItem("username");
        },
        nulldata() {
            this.message = null;
            this.switcher = null;
            this.conndata.hashname = null;
            this.conndata.address = null;
            this.conndata.public_key = null;
            this.userdata.username = null;
            this.userdata.password = null;
            this.userdata.password_repeat = null;
            this.userdata.private_key = null;
        },
        setswitch(name) {
            this.switcher = (this.switcher === name) ? null : name;
        },
    },
    mounted() {
        let token = localStorage.getItem("token");
        if (token) {
            this.authdata.token = token;
            this.authdata.username = localStorage.getItem("username");
        }
        this.opened = this.$route.name;
        switch (this.opened) {
            case RoutesData[routes.account].name: this.account(); break;
            case RoutesData[routes.network].name: this.network(); break;
            case RoutesData[routes.client].name: this.client(this.$route.params.id); break;
        }
    },
    updated() {
        if (this.opened == RoutesData[routes.network].name) {
            this.$nextTick(() => {
                var bottomChat = this.$refs.bottomChat;
                bottomChat.scrollTop = bottomChat.scrollHeight;
            });
        }
    },
    watch: {
        '$route' (to, from) {
            this.nulldata();
            this.opened = to.name;
        },
    },
});
