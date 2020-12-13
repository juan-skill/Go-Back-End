<template>
  <v-container>
    <v-container grid-list-xl>
      <v-layout>
        <v-flex>
          <v-card class="mb-3 pa-3">
            <v-btn column md-9 outlined color="indigo" @click="getDomains()">
              Show last domains
            </v-btn>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>

    <p v-if="loading">Loading...</p>
    <v-card v-if="showInfo" class="mx-auto">
      <v-list>
        <v-list-group prepend-icon="mdi-account-circle">
          <template v-slot:activator>
            <v-list-item-content>
              <v-list-item-title>Domains</v-list-item-title>
            </v-list-item-content>
          </template>

          <v-list-group
            no-action
            sub-group
            v-for="domain in domains"
            :key="domain.ssl_grade"
            prepend-icon="mdi-account-circle"
          >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Domain</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-group
              no-action
              sub-group
              v-for="(value, name) in domain"
              :key="value.tile"
            >
              <template v-slot:activator>
                <v-list-item-content>
                  <v-list-item-title>Domain {{ name }} </v-list-item-title>
                </v-list-item-content>
              </template>
              <v-list-item v-for="(value, name) in value" :key="value.title">
                <v-list-item-title
                  v-if="name != 'servers'"
                  v-text="name"
                ></v-list-item-title>

                <v-list-item-subtitle
                  v-if="name != 'servers'"
                  v-text="value"
                ></v-list-item-subtitle>
              </v-list-item>
            </v-list-group>
          </v-list-group>
        </v-list-group>
      </v-list>
    </v-card>
  </v-container>
</template>

<script>
import { mapActions, mapState } from "vuex";
export default {
  data() {
    return {
      domainName: ""
    };
  },
  computed: {
    ...mapState(["domains", "loading", "showInfo"])
  },
  methods: {
    ...mapActions(["getDomains"])
  }
};
</script>
