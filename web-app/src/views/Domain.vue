<template>
  <v-container>
    <v-container grid-list-xl>
      <v-layout>
        <v-flex>
          <v-card class="mb-3 pa-3">
            <form md-9 row @submit.prevent="getDomain(domainName)">
              <v-text-field
                column
                md-6
                label="Domain name"
                outlined
                v-model="domainName"
              ></v-text-field>
              <v-btn column md-9 outlined color="indigo" type="submit">
                Submit
              </v-btn>
            </form>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>

    <p v-if="submitting">Submitting...</p>
    <v-card class="mx-auto">
      <v-list>
        <v-list-group
          v-for="domain in domains"
          :key="domain.domainName"
          :value="true"
          prepend-icon="mdi-account-circle"
        >
          <template v-slot:activator>
            <v-list-item-title v-text="domain.domainName"></v-list-item-title>
          </template>

          <v-list-group :value="true" no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Title</v-list-item-title>
              </v-list-item-content>
            </template>
            <v-list-item link>
              <v-list-item-title v-text="domain.title"></v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Server Changed</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="domain in domains"
              :key="domain.domainName"
              link
            >
              <v-list-item-title
                v-text="domain.servers_changed"
              ></v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>SSL Grade</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="domain in domains"
              :key="domain.domainName"
              link
            >
              <v-list-item-title v-text="domain.ssl_grade"></v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Previous SSL Grade</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item link>
              <v-list-item-title
                v-text="domain.previous_ssl_grade"
              ></v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Logo</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="domain in domains"
              :key="domain.domainName"
              link
            >
              <v-list-item-title v-text="domain.logo"></v-list-item-title>
            </v-list-item>
          </v-list-group>

          <v-list-group no-action sub-group>
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Is Down?</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item
              v-for="domain in domains"
              :key="domain.domainName"
              link
            >
              <v-list-item-title v-text="domain.is_down"></v-list-item-title>
            </v-list-item>
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
    ...mapState(["domains", "submitting"])
  },
  methods: {
    ...mapActions(["getDomain"])
  }
};
</script>
