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
    <v-card v-if="showInfo" class="mx-auto">
      <v-list>
        <v-list-group prepend-icon="mdi-account-circle">
          <template v-slot:activator>
            <v-list-item-content>
              <v-list-item-title>Domain</v-list-item-title>
            </v-list-item-content>
          </template>

          <v-list-item
            v-for="(value, name) in domain"
            :key="name.ssl_grade"
            link
          >
            <v-list-item-title
              v-if="name != 'servers'"
              v-text="name"
            ></v-list-item-title>
            <v-list-item-subtitle
              v-if="name != 'servers'"
              v-text="value"
            ></v-list-item-subtitle>
          </v-list-item>
          <!-- v-if="name != 'servers'" -->

          <!-- <template v-slot:activator>
            <v-list-item-content>
              <v-list-item-title>{{ name }}</v-list-item-title>
            </v-list-item-content>
          </template>
          <v-list-item link>
            <v-list-item-title v-text="value"></v-list-item-title>
        
          </v-list-item>-->
        </v-list-group>
      </v-list>

      <v-list>
        <v-list-group no-action sub-group>
          <template v-slot:activator>
            <v-list-item-content>
              <v-list-item-title>servers</v-list-item-title>
            </v-list-item-content>
          </template>

          <v-list-group
            v-for="(val, nam) in domain.servers"
            :key="nam.address"
            prepend-icon="mdi-account-circle"
          >
            <template v-slot:activator>
              <v-list-item-content>
                <v-list-item-title>Server {{ nam }}</v-list-item-title>
              </v-list-item-content>
            </template>

            <v-list-item v-for="(val, nam) in val" :key="nam.address" link>
              <v-list-item-title v-text="nam"></v-list-item-title>
              <v-list-item-subtitle v-text="val"></v-list-item-subtitle>
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
    ...mapState(["domain", "submitting", "showInfo"])
  },
  methods: {
    ...mapActions(["getDomain"])
  }
};
</script>
