<template>
  <div>
    <v-app-bar app color="deep-purple accent-4" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Info Server</v-toolbar-title>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" absolute bottom temporary dark>
      <v-layout mt-4 column align-center>
        <v-list nav dense>
          <v-list-item-group
            v-model="group"
            active-class="deep-purple--text text--accent-4"
          >
            <v-list-item>
              <v-list-item-title>
                <RouterLink :to="{ name: 'Domain' }">Test domain</RouterLink>
              </v-list-item-title>
            </v-list-item>

            <v-list-item>
              <v-list-item-title>
                <RouterLink :to="{ name: 'Domains' }">Last records</RouterLink>
              </v-list-item-title>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </v-layout>
    </v-navigation-drawer>
  </div>
</template>

<script>
export default {
  name: "Header",
  data: () => ({
    drawer: false,
    group: null
  }),

  watch: {
    group() {
      this.drawer = false;
    }
  }
};
</script>
