<template>
  <div class="card floating">
    <div class="card-title">
      <h2>{{ t("prompts.newFile") }}</h2>
    </div>

    <div class="card-content">
      <span v-if="onlyOffice.url">
        <div
          v-for="option in this.typeOptions"
          class="card filetype-option file-icons"
          :class="{selected: selectedTypeOption == option, [option.name]: true}"
          @click="selectedTypeOption = option">
          <span :data-type=option.dataType :aria-label=option.template>
            <i class="material-icons"></i>
            <span class="name">{{option.name}}</span>
          </span>
        </div>
      </span>
      <p>{{ $t("prompts.newFileMessage") }}</p>
      <div class="suffixed-input">
        <input
          id="focus-prompt"
          class="input input--block"
          type="text"
          @keyup.enter="submit"
          v-model.trim="name"
        />
        <div v-if="selectedTypeOption.extension" class="suffix" :class="selectedTypeOption.name">
          {{selectedTypeOption.extension}}
        </div>
      </div>
    </div>

    <div class="card-action">
      <button
        class="button button--flat button--grey"
        @click="layoutStore.closeHovers"
        :aria-label="t('buttons.cancel')"
        :title="t('buttons.cancel')"
      >
        {{ t("buttons.cancel") }}
      </button>
      <button
        class="button button--flat"
        @click="submit"
        :aria-label="t('buttons.create')"
        :title="t('buttons.create')"
      >
        {{ t("buttons.create") }}
      </button>
    </div>
  </div>
</template>

<style scoped>
  .filetype-option {
    padding: 5px;
    display: inline-block;
    text-align: center;
    margin-right: 15px;
    font-size: 0.75rem;
    cursor: pointer;
  }

  .filetype-option i {
    font-size: 3rem;
  }

  .filetype-option span {
    display: block;
  }

  .filetype-option.selected, .filetype-option.selected i {
    color: white !important;
  }

  .filetype-option.selected.empty { background: rgba(0,0,0,0.8); }
  .filetype-option.selected.doc { background: var(--icon-blue); }
  .filetype-option.selected.sheet { background: var(--icon-green); }
  .filetype-option.selected.slide { background: var(--icon-orange); }
  .suffix.doc { background: var(--icon-blue) !important;  }
  .suffix.sheet { background: var(--icon-green) !important; }
  .suffix.slide { background: var(--icon-orange) !important; }
  .suffix { color: white; }
</style>

<script setup lang="ts">
import { inject, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { useFileStore } from "@/stores/file";
import { useLayoutStore } from "@/stores/layout";

import { files as api } from "@/api";
import url from "@/utils/url";
import { onlyOffice } from "@/utils/constants";

const $showError = inject<IToastError>("$showError")!;

const fileStore = useFileStore();
const layoutStore = useLayoutStore();

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const name = ref<string>("");
const typeOptions = ref<NewFileTemplate[]>([
  { name: "empty", extension: "", dataType: "text" },
  { name: "doc", extension: ".docx", template: "empty.docx" },
  { name: "sheet", extension: ".xlsx", template: "empty.xlsx" },
  { name: "slide", extension: ".pptx", template: "empty.pptx" },
]);
const selectedTypeOption = ref<NewFileTemplate>(typeOptions.value[0]);

const submit = async (event: Event) => {
  event.preventDefault();
  if (name.value === "") return;

  // Build the path of the new directory.
  let uri = fileStore.isFiles ? route.path + "/" : "/";

  if (!fileStore.isListing) {
    uri = url.removeLastDir(uri) + "/";
  }

  const filename = name.value + selectedTypeOption.value.extension;
  uri += encodeURIComponent(filename);
  uri = uri.replace("//", "/");

  if (selectedTypeOption.value.template) {
    uri = `${uri}?template=${selectedTypeOption.value.template}`;
  }

  try {
    await api.post(uri);
    router.push({ path: uri });
  } catch (e) {
    if (e instanceof Error) {
      $showError(e);
    }
  }

  layoutStore.closeHovers();
};
</script>
