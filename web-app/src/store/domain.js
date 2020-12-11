import Vue from "vue";
import Vuex from "vuex";

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    domain: null
  },
  mutations: {
    setDomain(state, payload) {
      state.domain = payload;
    }
  },
  actions: {
    getDomain({ commit }, domain) {
      //vaildat datos
      console.info(domain);
    }
  },
  modules: {}
});
