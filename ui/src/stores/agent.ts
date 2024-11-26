/**
 * AI Agent session management store
 */
import { ref } from 'vue'
import { defineStore } from 'pinia'
import { useAxiosClient } from '@/stores/axios'
const axios = useAxiosClient()

export type AgentMessage = {
  id: number
  author_id: number
  author_role: string
  agent_session_id: number
  content: string
  created_at: string
}

export type AgentSession = {
  id?: number
  description: string
  messages?: AgentMessage[]
  created_at?: string
}

export const useAgentStore = defineStore('agent', () => {
  const sessions = ref<AgentSession[]>([])
  const selectedSession = ref<AgentSession | null>(null)
  const errors = ref<string>('')


  const createSession = (session: AgentSession) => {
    return new Promise((resolve, reject) => {
      axios
        .post('/api/agent/sessions', session)
        .then((response) => {
          sessions.value.unshift(response.data)
          errors.value = '' // clear any previous errors
          resolve(response.data)
        })
        .catch((error) => {
          errors.value = error?.response?.data?.error ?? error.message
          reject(error)
        })
    })
  }

  /**
   * Get all ai agent sessions
   */
  const getSessions = () => {
    axios
      .get('/api/agent/sessions')
      .then((response) => {
        console.log("[agent.ts]", response.data)
        sessions.value = response.data
      })
      .catch((error) => {
        console.error("[agent.ts]", error)
      })
  }

  /**
   * Delete an ai agent session
   * @param session - The session to delete
   */
  const deleteSession = (session: AgentSession) => {
    axios
      .delete(`/api/agent/sessions/${session.id}`)
      .then(() => {
        sessions.value = sessions.value.filter((s) => s.id !== session.id)
      })
      .catch((error) => {
        console.error("[agent.ts]", error)
      })
  }

  /**
   * Send a user message to the ai agent session
   * @param message - The message to send
   */
  const sendUserMessage = (sessionId: number, message: string) => {
    axios
      .post(`/api/agent/sessions/${sessionId}/messages`, {
        content: message,
      })
      .then((response) => {
        console.info('[agent.ts] sendUserMessage', response.data)
      })
      .catch((error) => {
        console.error('[agent.ts] sendUserMessage', error)
      })
  }

  /**
   * Select an ai agent session
   * @param session - The session to select
   */
  const selectSession = (session: AgentSession) => {
    selectedSession.value = session
  }

  /**
   * Append a message to an ai agent session
   * @param session - The session to append the message to
   * @param message - The message to append
   */
  const appendMessageToSession = (sessionId: number, message: AgentMessage) => {
    const sessionToUpdate = sessions.value.find((s) => s.id === sessionId)
    if (sessionToUpdate) {
      if (!sessionToUpdate.messages) {
        sessionToUpdate.messages = []
      }
      sessionToUpdate.messages.push(message)
    } else {
      console.error('[agent.ts] appendMessageToSession: session not found, id: ', sessionId)
    }
  }

  return {
    createSession,
    getSessions,
    deleteSession,
    selectSession,
    errors,
    sessions,
    selectedSession,
    appendMessageToSession,
    sendUserMessage,
  }
})
