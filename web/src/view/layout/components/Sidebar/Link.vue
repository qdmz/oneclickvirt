<template>
  <component
    :is="linkProps(to).is"
    v-bind="linkProps(to).props"
    @click="handleNavigation"
  >
    <slot />
  </component>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { isExternal } from '@/utils/validate'

const props = defineProps({
  to: {
    type: String,
    required: true
  }
})

const router = useRouter()

const linkProps = (url) => {
  if (isExternal(url)) {
    return {
      is: 'a',
      props: {
        href: url,
        target: '_blank',
        rel: 'noopener'
      }
    }
  }
  return {
    is: 'div', // 改为div，避免router-link的自动导航
    props: {
      style: 'cursor: pointer;'
    }
  }
}

// 处理导航事件，确保立即跳转
const handleNavigation = (event) => {
  if (!isExternal(props.to)) {
    event.preventDefault()
    event.stopPropagation()
    
    // 立即进行路由跳转
    router.push(props.to)
  }
}
</script>