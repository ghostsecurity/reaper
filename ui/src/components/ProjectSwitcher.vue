<template>
  <Select v-model="selectedProject">
    <SelectTrigger aria-label="Select account"
                   :class="cn(
                    'flex items-center gap-2 [&>span]:line-clamp-1 [&>span]:flex [&>span]:w-full [&>span]:items-center [&>span]:gap-1 [&>span]:truncate [&_svg]:h-4 [&_svg]:w-4 [&_svg]:shrink-0',
                    { 'flex h-9 w-9 shrink-0 items-center justify-center p-0 [&>span]:w-auto [&>svg]:hidden': isCollapsed },
                  )">
      <SelectValue placeholder="Select an account">
        <div class="flex items-center gap-3">
          <Icon class="size-4"
                :icon="selectedProjectData!.icon" />
          <span v-if="!isCollapsed">
            {{ selectedProjectData!.label }}
          </span>
        </div>
      </SelectValue>
    </SelectTrigger>
    <SelectContent>
      <SelectItem v-for="project of projects"
                  :key="project.value"
                  :value="project.value">
        <div class="flex items-center gap-3 [&_svg]:size-4 [&_svg]:shrink-0 [&_svg]:text-foreground">
          <Icon class="size-4"
                :icon="project.icon" />
          {{ project.label }}
        </div>
      </SelectItem>
    </SelectContent>
  </Select>
</template>


<script lang="ts" setup>
import { computed, ref } from 'vue'
import { Icon } from '@iconify/vue'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { cn } from '@/lib/utils'

interface ProjectSwitcherProps {
  isCollapsed: boolean
  projects: {
    label: string
    value: string
    icon: string
  }[]
}

const props = defineProps<ProjectSwitcherProps>()

const selectedProject = ref<string>(props.projects[0].value)
const selectedProjectData = computed(() => props.projects.find(item => item.value === selectedProject.value))
</script>
