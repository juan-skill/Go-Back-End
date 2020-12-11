<template>
  <v-card>
    <v-container grid-list-xl>
      <v-layout>
        <v-flex>
          <v-card class="mb-3 pa-3">
            <form md-9 row @submit.prevent="fetchUsers">
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

    <v-row class="pa-4" justify="space-between">
      <v-col cols="5">
        <v-treeview
          :active.sync="active"
          :items="items"
          :load-children="fetchUsers"
          :open.sync="open"
          activatable
          color="warning"
          open-on-click
          transition
        >
          <template v-slot:prepend="{ item }">
            <v-icon v-if="!item.children">
              mdi-account
            </v-icon>
          </template>
        </v-treeview>
      </v-col>

      <v-divider vertical></v-divider>

      <v-col class="d-flex text-center">
        <v-scroll-y-transition mode="out-in">
          <div
            v-if="!selected"
            class="title grey--text text--lighten-1 font-weight-light"
            style="align-self: center;"
          >
            Select a User
          </div>
          <v-card
            v-else
            :key="selected.domainName"
            class="pt-6 mx-auto"
            flat
            max-width="400"
          >
            <v-card-text>
              <v-avatar v-if="avatar" size="88">
                <v-img
                  :src="`https://avataaars.io/${avatar}`"
                  class="mb-6"
                ></v-img>
              </v-avatar>
              <h3 class="headline mb-2">
                {{ selected.grade_ssl }}
              </h3>
              <div class="blue--text mb-2">
                {{ selected.email }}
              </div>
              <div class="blue--text subheading font-weight-bold">
                {{ selected.username }}
              </div>
            </v-card-text>
            <v-divider></v-divider>
            <v-row class="text-left" tag="v-card-text">
              <v-col class="text-right mr-4 mb-2" tag="strong" cols="5">
                Company:
              </v-col>
              <v-col>{{ selected.company.name }}</v-col>
              <v-col class="text-right mr-4 mb-2" tag="strong" cols="5">
                Website:
              </v-col>
              <v-col>
                <a :href="`//${selected.website}`" target="_blank">{{
                  selected.website
                }}</a>
              </v-col>
              <v-col class="text-right mr-4 mb-2" tag="strong" cols="5">
                Phone:
              </v-col>
              <v-col>{{ selected.phone }}</v-col>
            </v-row>
          </v-card>
        </v-scroll-y-transition>
      </v-col>
    </v-row>
  </v-card>
  <!-- <div class="domain">
    <h1>Info Server</h1>
    <form @submit.prevent="showDomain()">
      <input type="text" placeholder="domain" v-model="domainName" />
      <button type="submit">Submit</button>
      <p>Domain is: {{ domainName }}</p>
    </form>
    <div>
      <v-container>
        <v-row>
          <v-col> {{ domain }} </v-col>
        </v-row>
      </v-container>
    </div>
  </div>
  -->
</template>

<script>
// @ is an alias to /src
const avatars = [
  "?accessoriesType=Blank&avatarStyle=Circle&clotheColor=PastelGreen&clotheType=ShirtScoopNeck&eyeType=Wink&eyebrowType=UnibrowNatural&facialHairColor=Black&facialHairType=MoustacheMagnum&hairColor=Platinum&mouthType=Concerned&skinColor=Tanned&topType=Turban",
  "?accessoriesType=Sunglasses&avatarStyle=Circle&clotheColor=Gray02&clotheType=ShirtScoopNeck&eyeType=EyeRoll&eyebrowType=RaisedExcited&facialHairColor=Red&facialHairType=BeardMagestic&hairColor=Red&hatColor=White&mouthType=Twinkle&skinColor=DarkBrown&topType=LongHairBun",
  "?accessoriesType=Prescription02&avatarStyle=Circle&clotheColor=Black&clotheType=ShirtVNeck&eyeType=Surprised&eyebrowType=Angry&facialHairColor=Blonde&facialHairType=Blank&hairColor=Blonde&hatColor=PastelOrange&mouthType=Smile&skinColor=Black&topType=LongHairNotTooLong",
  "?accessoriesType=Round&avatarStyle=Circle&clotheColor=PastelOrange&clotheType=Overall&eyeType=Close&eyebrowType=AngryNatural&facialHairColor=Blonde&facialHairType=Blank&graphicType=Pizza&hairColor=Black&hatColor=PastelBlue&mouthType=Serious&skinColor=Light&topType=LongHairBigHair",
  "?accessoriesType=Kurt&avatarStyle=Circle&clotheColor=Gray01&clotheType=BlazerShirt&eyeType=Surprised&eyebrowType=Default&facialHairColor=Red&facialHairType=Blank&graphicType=Selena&hairColor=Red&hatColor=Blue02&mouthType=Twinkle&skinColor=Pale&topType=LongHairCurly"
];

const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

export default {
  data: () => ({
    domainName: "",
    active: [],
    avatar: null,
    open: [],
    users: []
  }),
  computed: {
    items() {
      return [
        {
          name: "Users",
          children: this.users
        }
      ];
    },
    selected() {
      if (!this.active.length) return undefined;

      const id = this.active[0];

      return this.users.find(user => user.id === id);
    }
  },

  watch: {
    selected: "randomAvatar"
  },

  methods: {
    async fetchUsers(item) {
      console.log(this.domainName);

      try {
        await pause(1500);

        const resp = await fetch("http://localhost:8090/get-last-domains", {
          method: "GET"
        });

        const data = await resp.json();
        console.log(data);

        item.children.push(...data);
      } catch (error) {
        console.warn(error);
      }
    },
    /*
    async fetchUsers (item) {
      // Remove in 6 months and say
      // you've made optimizations! :)
      await pause(1500)

      return fetch('https://jsonplaceholder.typicode.com/users')
        .then(
          res => res.json()
        )
        .then(json => (item.children.push(...json)))
        .catch(err => console.warn(err))
    },
    */
    randomAvatar() {
      this.avatar = avatars[Math.floor(Math.random() * avatars.length)];
      console.info(this.avatar);
    }
  }
  /*
  name: "domain",
  data() {
    return {
      domainName: "",
      domain: null
    };
  },
  methods: {
    async showDomain() {
      console.log(this.domainName);
      try {
        const resp = await fetch("http://localhost:8090/domain", {
          method: "POST",
          headers: {
            "Content-Type": "application/json"
          },
          body: JSON.stringify({ domainName: this.domainName })
        });

        const resDB = await resp.json();
        console.log(resDB);
        this.domain = resDB;
        console.info(this.domain.logo);
      } catch (error) {
        console.error(error);
      }
    }
  }
  */
};
</script>
