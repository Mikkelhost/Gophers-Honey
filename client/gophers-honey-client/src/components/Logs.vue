<template>
  <div style="overflow-x: hidden; overflow-y: hidden; height: 100%">
    <b-row>
      <b-col lg="6" class="my-1">
        <b-form-group
            label="Filter"
            label-for="filter-input"
            label-cols-sm="3"
            label-align-sm="right"
            label-size="sm"
            class="mb-0"
        >
          <b-input-group size="sm">
            <b-form-input
                id="filter-input"
                v-model="filter"
                type="search"
                placeholder="Type to Search"
            ></b-form-input>

            <b-input-group-append>
              <b-button :disabled="!filter" @click="filter = ''">Clear</b-button>
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
      </b-col>
      <b-col lg="8" class="my-1">
        <b-form-group
            v-model="sortDirection"
            label="Filter On"
            description="Leave all unchecked to filter on all data"
            label-cols-sm="3"
            label-align-sm="right"
            label-size="sm"
            class="mb-0"
            v-slot="{ ariaDescribedby }"
        >
          <b-form-checkbox-group
              v-model="filterOn"
              :aria-describedby="ariaDescribedby"
              class="mt-1"
          >
            <b-form-checkbox value="device_id">Device id</b-form-checkbox>
            <b-form-checkbox value="src_host">Source Host</b-form-checkbox>
            <b-form-checkbox value="src_port">Source Port</b-form-checkbox>
            <b-form-checkbox value="dst_host">Destination Host</b-form-checkbox>
            <b-form-checkbox value="dst_port">Destination Port</b-form-checkbox>
          </b-form-checkbox-group>
        </b-form-group>
      </b-col>
    </b-row>
    <b-table bordered sticky-header striped hover :items="logs" :fields="fields"
             style="height: 735px; max-height: 735px!important;"
             :filter="filter"
             :filter-included-fields="filterOn"
    >
      <template #cell(show_details)="row">
        <b-button size="md" @click="row.toggleDetails" class="mr-2">
          {{ row.detailsShowing ? 'Hide' : 'Show' }} Raw Logs
        </b-button>
      </template>

      <template #row-details="row">
        <b-card>
          <b-row class="mb-2">
            <b-col>{{ row.item.raw_log }}</b-col>
          </b-row>
        </b-card>
      </template>
    </b-table>
  </div>
</template>

<script>

export default {
  name: "logs",
  props: {
    logs: {
      type: Array,
    }
  },
  mounted() {
    console.log(this.logs)
  },
  methods: {
  },
  data: function () {
    return {
      filter: null,
      filterOn: [],
      sortDirection: 'asc',
      fields: [
        {
          key: "device_id",
          sortable: true
        },
        {
          key: "level",
          label: "Log Level",
          sortable: true
        },
        {
          key: "log_type",
          sortable: true
        },
        {
          key: "src_host",
          label: "Source Host",
          sortable: true
        },
        {
          key: "src_port",
          label: "Source Port",
          sortable: true
        },
        {
          key: "dst_host",
          label: "Destination Host",
          sortable: true
        },
        {
          key: "dst_port",
          label: "Destination Port",
          sortable: true
        },
        {
          key: "message",
          label: "Log Data",
          sortable: false
        },
        {
          key: "log_time_stamp",
          label: "Timestamp(UTC)",
          sortable: true
        },
        {
          key: "show_details",
          label: "Show Raw Logs"
        }
      ]
    }
  }
}
</script>

<style scoped>

</style>