<template>
  <div class="flex h-full flex-col">
    <div class="flex items-center p-2">
      <div class="flex items-center gap-2">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost"
                    size="icon"
                    :disabled="!endpoint"
                    @click="clearResults">
              <Trash2 class="size-4" />
              <span class="sr-only">Move to trash</span>
            </Button>
          </TooltipTrigger>
          <TooltipContent>Clear all results</TooltipContent>
        </Tooltip>
      </div>
      <div class="ml-auto flex items-center gap-2">
        <Dialog v-model:open="isModalOpen">
          <DialogTrigger as-child>
            <Button variant="outline"
                    :disabled="!endpoint || attackRunning">
              <PocketKnifeIcon class="mr-2 size-4" />Create a test
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>New dynamic test</DialogTitle>
              <DialogDescription>
                Create and run a new dynamic test on this endpoint.
              </DialogDescription>
              <DialogDescription class="rounded-md bg-muted p-2 text-xs font-medium">
                <RequestMethod :code="-1">
                  {{ endpoint?.method }}
                </RequestMethod>
                <span class="ml-2 text-xs">
                  {{ endpoint?.path }}
                </span>
              </DialogDescription>
            </DialogHeader>
            <span class="text-xs font-medium">
              Test type
            </span>
            <Popover v-model:open="attackTemplateSelectOpen">
              <PopoverTrigger as-child>
                <Button variant="outline"
                        role="combobox"
                        :aria-expanded="attackTemplateSelectOpen"
                        class="w-full justify-between">
                  {{ attackTemplateSelectValue
                    ? testTypes.find((tt) => tt.value === attackTemplateSelectValue)?.label
                    : "Select test type..." }}
                  <CaretSortIcon class="ml-2 size-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent class="w-full p-0">
                <Command>
                  <CommandInput class="h-9"
                                placeholder="Search test type..." />
                  <CommandEmpty>No test type found.</CommandEmpty>
                  <CommandList>
                    <CommandGroup>
                      <CommandItem v-for="tt in testTypes"
                                   :key="tt.value"
                                   :value="tt.value"
                                   @select="(ev) => {
                                    if (typeof ev.detail.value === 'string') {
                                      attackTemplateSelectValue = ev.detail.value
                                    }
                                    attackTemplateSelectOpen = false
                                  }">
                        {{ tt.label }}
                        <CheckIcon class="ml-auto size-4"
                                   :class="attackTemplateSelectValue === tt.value ? 'opacity-100' : 'opacity-0'" />
                      </CommandItem>
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
            <div v-if="attackTemplateSelectValue === 'idor'"
                 class="flex flex-col gap-4">
              <span class="text-xs font-medium">
                Included parameters
              </span>
              <div v-if="endpoint"
                   class="mx-1 space-y-2 text-muted-foreground">
                <div v-for="param in endpoint.params.split(',')"
                     :key="param"
                     class="flex items-center gap-2">
                  <Checkbox :checked="endpointStore.isParamSelected(param)"
                            @update:checked="(checked) => endpointStore.toggleParam(param, checked)" />
                  <label
                         class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                    {{ param }}
                  </label>
                </div>
              </div>
              <span class="text-xs text-muted-foreground">These parameters were dynamically determined from previous
                requests to this endpoint. Included parameters will be used as test inputs, while the rest of parameters
                will be left intact from the original request.</span>
            </div>
            <div v-if="unsupportedTest">
              <span class="text-xs font-medium">
                Unsupported test type
              </span>
              <div class="mt-1 space-y-2 text-xs text-muted-foreground">
                We'll be adding more test types soon. In the meantime, please submit an issue on <a
                   href="https://github.com/ghostsecurity/reaper/issues"
                   target="_blank"
                   class="text-foreground">GitHub</a>.
              </div>
            </div>
            <DialogFooter>
              <div>
                <Button type="submit"
                        @click.prevent="handleCreateTest"
                        :disabled="attackRunning"
                        v-if="!attackRunning && !attackComplete && !unsupportedTest">
                  <PocketKnifeIcon class="mr-2 size-4" />
                  {{ attackRunning ? 'Starting...' : 'Start test' }}
                </Button>
                <Button class="w-full"
                        :disabled="attackRunning"
                        v-if="attackRunning">
                  <RefreshCwIcon class="mr-2 size-4 animate-spin" /> Test running...
                </Button>
                <Button class="w-full"
                        v-if="attackComplete && !attackRunning"
                        @click="isModalOpen = false">
                  <ScrollText class="mr-2 size-4" /> View results
                </Button>
              </div>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
    <Separator />
    <div v-if="endpoint"
         class="flex flex-1 flex-col">
      <div class="flex items-start p-4">
        <div class="mr-4 flex items-start gap-4 text-sm">
          <div class="flex flex-col gap-2">
            <RequestMethod :code="-1">
              {{ endpoint.method }}
            </RequestMethod>
          </div>
          <div class="grid gap-1">
            <div class="mr-4 truncate font-semibold text-foreground/80">
              {{ endpoint.path }}
            </div>
            <div class="space-y-1 text-xs">
              <div class="font-medium text-muted-foreground">
                {{ endpoint.hostname }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="endpoint.created_at"
             class="ml-auto space-y-4 text-xs text-muted-foreground">
          <div>
            {{ format(new Date(endpoint.created_at), "PPpp") }}
          </div>
        </div>
      </div>
      <Separator />
      <div class="h-screen justify-center space-y-2 overflow-y-auto bg-muted/50 p-1 text-sm">
        <div v-for="result in endpointStore.results"
             :key="result.id"
             class="flex flex-col items-start gap-2 rounded-md bg-background/95 p-3 text-left text-sm shadow-sm transition-all hover:bg-accent/50">
          <Badge variant="outline"
                 class="border-green-600/20 bg-green-50 px-1 text-green-700">
            <div class="text-2xs font-semibold uppercase">
              Success
            </div>
          </Badge>
          <span class="text-xs text-muted-foreground">{{ result.request }}</span>
        </div>
        <div v-if="endpointStore.results.length < 1"
             class="p-8 text-center font-medium text-muted-foreground">
          No results yet. Start a dynamic test.
        </div>
      </div>
    </div>
    <div v-else
         class="p-8 text-center text-muted-foreground">
      No endpoint selected.
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { PocketKnifeIcon, RefreshCwIcon, ScrollText, Trash2 } from 'lucide-vue-next'
import { format } from 'date-fns'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Badge } from '@/components/ui/badge'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import { CaretSortIcon, CheckIcon } from '@radix-icons/vue'
import type { Endpoint } from '@/stores/endpoint'
import { useEndpointStore } from '@/stores/endpoint'
import RequestMethod from './shared/RequestMethod.vue'

interface EndpointDisplayProps {
  endpoint: Endpoint | undefined
}

const endpointStore = useEndpointStore()
const attackTemplateSelectOpen = ref(false)
const attackTemplateSelectValue = ref('')
const attackRunning = computed(() => endpointStore.attackRunning)
const attackComplete = computed(() => endpointStore.attackComplete)
const props = defineProps<EndpointDisplayProps>()

const isModalOpen = ref(false)

const handleCreateTest = () => {
  if (props.endpoint) {
    endpointStore.startAttack(props.endpoint)
  }
}

const clearResults = () => {
  if (props.endpoint) {
    endpointStore.clearResults(props.endpoint)
  }
}

const testTypes = [
  { value: 'idor', label: 'Insecure Direct Object Reference (IDOR/BOLA)' },
  { value: 'bf', label: 'Brute Force' },
  { value: 'sqli', label: 'SQL Injection (SQLi)' },
  { value: 'header', label: 'HTTP Header Injection' },
  { value: 'xml', label: 'XML Injection (XXE)' },
  { value: 'xss', label: 'Cross Site Scripting (XSS)' },
  { value: 'ssrf', label: 'Server Side Request Forgery (SSRF)' },
  { value: 'lfi', label: 'Local File Inclusion (LFI)' },
  { value: 'patht', label: 'Path Traversal' },
  { value: 'ssti', label: 'Server Side Template Injection (SSTI)' },
  { value: 'cmdi', label: 'Command Injection (cmd/shell)' },
]

const unsupportedTest = computed(() => {
  return attackTemplateSelectValue.value !== 'idor' && attackTemplateSelectValue.value.length > 0
})
</script>
