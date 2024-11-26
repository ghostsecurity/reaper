<template>
  <Table>
    <TableCaption class="text-xs">Recently scanned domains.</TableCaption>
    <TableHeader class="text-xs">
      <TableRow>
        <TableHead class="w-[240px]">
          Domain
        </TableHead>
        <TableHead class="text-center">
          Hosts
        </TableHead>
        <TableHead class="text-center">
          Last Scan
        </TableHead>
      </TableRow>
    </TableHeader>
    <TableBody class="text-xs">
      <TableRow v-for="(domain) in filteredDomains"
                :key="domain.name"
                class="cursor-pointer"
                @click="handleRowClick(domain)"
                :class="selectedDomain?.id === domain.id ? 'bg-muted/50' : ''">
        <TableCell class="flex items-center font-normal">
          <span class="mr-1 w-0.5 pr-0.5"
                :class="selectedDomain?.id === domain.id ? 'bg-primary' : 'bg-transparent'">&nbsp;</span>
          {{ domain.name }} <span class="ml-2 size-4">
            <span v-if="domain.status !== 'completed'"
                  class="rounded-sm bg-secondary px-1.5 py-0.5 text-muted-foreground">{{ domain.status }}</span>
          </span>
        </TableCell>
        <TableCell class="text-center font-medium">{{ domain.host_count && domain.host_count > 0 &&
          utils.customNumberFormat(domain.host_count) }}
        </TableCell>
        <TableCell class="text-center text-muted-foreground">
          {{ domain.last_scanned_at && utils.timeAgoInWords(domain.last_scanned_at) }}
        </TableCell>
      </TableRow>
    </TableBody>
  </Table>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { useScanStore } from '@/stores/scan'
import { useUtilStore } from '@/utils'
import type { Domain } from '@/stores/scan'

const scanStore = useScanStore()
const domains = computed(() => scanStore.domains)
const selectedDomain = ref<Domain | null>(null)

interface ScanListProps {
  filter: string
}

const utils = useUtilStore()
const props = defineProps<ScanListProps>()

const filteredDomains = computed(() => {
  if (props.filter === '') {
    return domains.value
  }
  return domains.value.filter((domain) => domain.name.includes(props.filter))
})

/**
 * Handle keyboard navigation
 * @param event - The keyboard event
 */
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'ArrowUp') {
    event.preventDefault()
    moveSelection(-1)
  } else if (event.key === 'ArrowDown') {
    event.preventDefault()
    moveSelection(1)
  }
}

/**
 * Move the selection up or down
 * @param direction - The direction to move the selection
 */
const moveSelection = (direction: number) => {
  const newIndex = selectedDomain.value ? filteredDomains.value.indexOf(selectedDomain.value) + direction : 0
  if (newIndex >= 0 && newIndex < filteredDomains.value.length) {
    selectedDomain.value = filteredDomains.value[newIndex]
    scanStore.selectDomain(filteredDomains.value[newIndex])
  }
}

const handleRowClick = (domain: Domain) => {
  selectedDomain.value = domain
  scanStore.selectDomain(domain)
}

onMounted(() => {
  scanStore.getDomains()

  // Add event listener for keyboard navigation
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  // Remove event listener when component is unmounted
  window.removeEventListener('keydown', handleKeyDown)
})
</script>
