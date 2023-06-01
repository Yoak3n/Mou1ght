import { defineStore } from 'pinia'
import { ss } from '@/utils/storage'

const LOCAL_NAME = 'appSetting'

export type Theme = 'light' | 'dark' | 'auto'

export type Language = 'zh-CN' | 'zh-TW' | 'en-US'

export interface AppState {
  siderCollapsed: boolean
  theme: Theme
  language: Language
}

export function defaultSetting(): AppState {
  return { siderCollapsed: false, theme: 'auto', language: 'en-US' }
}

export function getLocalSetting(): AppState {
  const localSetting: AppState | undefined = ss.get(LOCAL_NAME)
  return { ...defaultSetting(), ...localSetting }
}

export function setLocalSetting(setting: AppState): void {
  ss.set(LOCAL_NAME, setting)
}


export const useAppStore = defineStore('app-store', {
    state: (): AppState => getLocalSetting(),
    actions: {
      setSiderCollapsed(collapsed: boolean) {
        this.siderCollapsed = collapsed
        this.recordState()
      },
  
      setTheme(theme: Theme) {
        this.theme = theme
        this.recordState()
      },
  
      setLanguage(language: Language) {
        if (this.language !== language) {
          this.language = language
          this.recordState()
        }
      },
  
      recordState() {
        setLocalSetting(this.$state)
      },
    },
  })