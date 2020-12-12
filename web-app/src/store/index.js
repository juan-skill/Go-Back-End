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
  },
  mutations: {
    setDomain(state, payload) {
      state.domain = payload;
    },
    setDomains(state, payload) {
      state.domains.push(payload);
    },
    setResetDomains(state, payload) {
      state.domains = payload;
    },
    setSubmit(state, payload) {
      state.submitting = payload;
    },
    setLoading(state, payload) {
      state.loading = payload;
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
    },
    async getDomains({ commit }) {
      try {
        commit("setLoading", true);
        commit("setResetDomains", []);
        const pause = ms => new Promise(resolve => setTimeout(resolve, ms));
        await pause(400);

        const response = await axios.get(
          "http://localhost:8090/get-last-domains"
        );
        commit("setLoading", false);

        console.info(response.data);

        commit("setResetDomains", response.data);
      } catch (error) {
        console.warn(error);
      }
    }
  },
  modules: {}
});
