<template>
  <div>
    <v-app-bar app color="deep-purple accent-4" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Info Server</v-toolbar-title>
    </v-app-bar>

    <v-navigation-drawer v-model="drawer" absolute bottom temporary dark>
      <v-layout mt-4 column align-center>
        <v-list dense>
      <v-subheader>Links</v-subheader>
      <v-list-item-group>
        
        <v-list-item
          v-for="(item, i) in items"
          :key="i"
          router :to="item.route"
        >
                  <v-list-item-icon>
            <v-icon v-text="item.icon"></v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="item.text"></v-list-item-title>
          </v-list-item-content>
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
    group: null,
    items: [
        { text: 'Test domain', icon: 'mdi-clock', route: '/domain' },
        { text: 'Last records', icon: 'mdi-account', route: '/domains' },
      ],
  }),

  watch: {
    group() {
      this.drawer = false;
    }
  }
};
</script>
