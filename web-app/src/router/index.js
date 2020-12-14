import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Domain from "../views/Domain.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/domain",
    name: "Domain",
    component: Domain
  },
  {
    path: "/domains",
    name: "Domains",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/Domains.vue")
  }
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes
});

export default router;
