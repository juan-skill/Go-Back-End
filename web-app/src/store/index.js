import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    domain: {
      servers: [],
      servers_changed: false,
      ssl_grade: "",
      previous_ssl_grade: "",
      logo: "",
      title: "",
      is_down: false
    },
    loading: false,
    submitting: false,
    showInfo: false,
    domains: [],
    messageError: {
      message: "",
      active: false
    }
  },
  mutations: {
    setDomain(state, payload) {
      state.domain.ssl_grade = payload.ssl_grade;
      state.domain.previous_ssl_grade = payload.previous_ssl_grade;
      state.domain.logo = payload.logo;
      state.domain.title = payload.title;
      state.domain.is_down = payload.is_down;

      for (let i = 0; i < payload.servers.length; i++) {
        let server = new Object();
        server.address = payload.servers[i].address;
        server.country = payload.servers[i].country;
        server.owner = payload.servers[i].owner;
        server.ssl_grade = payload.servers[i].ssl_grade;
        state.domain.servers.push(server);
      }
    },
    setDomains(state, payload) {
      state.domains.push(payload);
    },
    setResetDomains(state, payload) {
      for (let i = 0; i < payload.length; i++) {
        let domain = new Object();
        domain.ssl_grade = payload[i].ssl_grade;
        domain.previous_ssl_grade = payload[i].previous_ssl_grade;
        domain.logo = payload[i].logo;
        domain.title = payload[i].title;
        domain.is_down = payload[i].is_down;
        state.domains.push(domain);
      }
    },
    setSubmit(state, payload) {
      state.submitting = payload;
    },
    setLoading(state, payload) {
      state.loading = payload;
    },
    setShowInfo(state, payload) {
      state.showInfo = payload;
    },
    setMessageError(state, payload) {
      state.messageError.message = payload.message;
      state.messageError.active = payload.active;
    }
  },
  actions: {
    async getDomain({ commit }, domainName) {
      commit("setSubmit", true);
      commit("setShowInfo", false);

      if (domainName  == "") {
        const err = {
          message: "cannot be empty domain name",
          active: true
        };
        commit("setMessageError", err);        
        commit("setSubmit", false);
      } else {

        try {

        const pause = ms => new Promise(resolve => setTimeout(resolve, ms));
        await pause(1500);

        const bodyRequest = { domainName: domainName };
        const headersRequest = { "Content-type": "application/json" };
        const response = await axios.post(
          "http://localhost:8090/domain",
          bodyRequest,
          { headersRequest }
        );

        console.info(response);
        

        if (response.statusText == "Created") {
          console.info(response.status);
          commit("setDomain", response.data);
          commit("setShowInfo", true);
        }
        commit("setSubmit", false);
      } catch (error) {
        commit("setSubmit", false);
        const err = {
          message: "try again in minutes",
          active: true
        };
        commit("setMessageError", err);
        console.warn(error);
      }
    }
    },
    async getDomains({ commit }) {
      try {
        commit("setLoading", true);
        commit("setShowInfo", false);
        commit("setResetDomains", []);
        const pause = ms => new Promise(resolve => setTimeout(resolve, ms));
        await pause(400);

        const response = await axios.get(
          "http://localhost:8090/get-last-domains"
        );
        commit("setLoading", false);

        console.info(response.data);
        console.info(response.statusText);
        //if (response.statusText == "OK") {
        if (response.statusText == "OK") {
          commit("setDomains", response.data);
          commit("setShowInfo", true);
        }
        commit("setLoading", false);
      } catch (error) {
        commit("setLoading", false);
        const err = {
          message: "try again in minutes",
          active: true
        };
        commit("setMessageError", err);
        console.warn(error);
      }
    }
  },
  modules: {}
});
