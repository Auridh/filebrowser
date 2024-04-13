<template>
  <div id="editor-container" :class="onlyOffice.fullViewport && 'full-viewport'">
    <header-bar v-if="!onlyOffice.fullViewport">
      <action icon="close" :label="$t('buttons.close')" @action="close()" />
      <title>{{ fileStore.req.name }}</title>
    </header-bar>
    <Breadcrumbs  v-if="!onlyOffice.fullViewport" base="/files" noLink />
    <errors v-if="error" :errorCode="error.status" />
    <div id="editor">
      <div id="onlyoffice-editor"></div>
    </div>
  </div>
</template>

<style scoped>
  #editor-container.full-viewport {
    padding-top: 0;
  }
</style>

<script setup lang="ts">
import url from "@/utils/url";
import { onlyOffice } from "@/utils/constants";

import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";
import Breadcrumbs from "@/components/Breadcrumbs.vue";
import Errors from "@/views/Errors.vue";
import { fetchJSON, StatusError } from "@/api/utils";
import { useFileStore } from "@/stores/file.js";
import { useRoute, useRouter } from "vue-router";
import { onBeforeUnmount, onMounted, ref } from "vue";

const fileStore = useFileStore();
const route = useRoute();
const router = useRouter();

const error = ref<StatusError | null>(null);
const editor = ref<DocsAPI.DocEditor | null>(null);

onMounted(() => {
  const isMobile = window.innerWidth <= 736;
  const clientConfigPromise = fetchJSON(
    `/api/onlyoffice/client-config${fileStore.req.path}?isMobile=${isMobile}`);
  window.addEventListener("keydown", keyEvent);

  const scriptUrl = `${onlyOffice.url}/web-apps/apps/api/documents/api.js`;
  const onlyofficeScript = document.createElement("script");
  onlyofficeScript.setAttribute("src", scriptUrl);
  document.head.appendChild(onlyofficeScript);

  onlyofficeScript.onload = async () => {
    try {
      const clientConfig = await clientConfigPromise;
      // eslint-disable-next-line no-undef
      editor.value = new DocsAPI.DocEditor("onlyoffice-editor", clientConfig);
    } catch (e) {
      error.value = e;
    }
  };
});

onBeforeUnmount(() => {
  window.removeEventListener("keydown", keyEvent);
  editor.value.destroyEditor();
});

const back = () => {
  let uri = url.removeLastDir(route.path) + "/";
  router.push({ path: uri });
};
const keyEvent = (event) => {
  if (!event.ctrlKey && !event.metaKey) {
    return;
  }

  if (String.fromCharCode(event.which).toLowerCase() !== "s") {
    return;
  }

  event.preventDefault();
  this.save();
};
const close = () => {
  fileStore.updateRequest(null)

  let uri = url.removeLastDir(route.path) + "/";
  router.push({ path: uri });
};
</script>
