import Vue from "vue";
import Vuex from "vuex";
import axios from "axios";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    domain: null,
    loading: false,
    submitting: false,
    domains: []
    /*
    servers: {
      address: "",
      country: "",
      owner: "",
      ssl_grade: ""
    },
    domain: {
      servers: [],
      servers_changed: "",
      ssl_grade: undefined,
      logo: "",
      title: "",
      is_down: undefined
    }
    */
  },
  mutations: {
    setDomain(state, payload) {
      state.domain = payload;
      //state.domain.ssl_grade = payload.ssl_grade;
      //console.info(state.domain.ssl_grade)
    },
    setDomains(state, payload) {
      state.domains.push(payload);
    },
    setResetDomains(state, payload) {
      state.domains = payload;
    },
    setSubmit(state, payload) {
      state.submitting = payload;
    }
  },
  actions: {
    async getDomain({ commit }, domainName) {
      //vaildat datos

      try {
        commit("setSubmit", true);
        commit("setResetDomains", []);
        const pause = ms => new Promise(resolve => setTimeout(resolve, ms));
        await pause(1500);

        const bodyRequest = { domainName: domainName };
        const headersRequest = { "Content-type": "application/json" };
        const response = await axios.post(
          "http://localhost:8090/domain",
          bodyRequest,
          { headersRequest }
        );
        commit("setSubmit", false);

        /*
        const resp = await fetch("http://localhost:8090/domain", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({ domainName: domainName })
        });

        const data = await resp.json();
        */
        console.info(response);
        /*
        const stringif = JSON.stringify(data);
        const parse = JSON.parse(stringif);
        console.warn(parse);
        */
        //commit("setDomain", response.data);
        commit("setDomains", response.data);
      } catch (error) {
        console.warn(error);
      }
    }
  },
  modules: {}
});
