import axios from 'axios'
import type {
  Bookmark,
  BookmarkFormData,
  CreateBookmarkInput,
  UpdateBookmarkInput,
  PaginatedResponse,
  SearchParams,
  TagInfo,
  HealthCheckResponse,
} from '../types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor for logging
api.interceptors.request.use(
  (config) => {
    console.log('API Request:', config.method?.toUpperCase(), config.url)
    return config
  },
  (error) => {
    console.error('API Request Error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Response Error:', error.response?.status, error.response?.data)
    return Promise.reject(error)
  }
)

// Bookmark API functions
export const createBookmark = async (data: CreateBookmarkInput): Promise<Bookmark> => {
  try {
    const response = await api.post<Bookmark>('/bookmarks', data)
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to create bookmark')
    }
    throw error
  }
}

export const getBookmark = async (id: string): Promise<Bookmark> => {
  try {
    const response = await api.get<Bookmark>(`/bookmarks/${id}`)
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to get bookmark')
    }
    throw error
  }
}

export const getBookmarks = async (
  page = 1,
  pageSize = 10
): Promise<PaginatedResponse<Bookmark>> => {
  try {
    const response = await api.get<PaginatedResponse<Bookmark>>('/bookmarks', {
      params: { page, pageSize },
    })
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to get bookmarks')
    }
    throw error
  }
}

export const updateBookmark = async (
  id: string,
  data: UpdateBookmarkInput
): Promise<Bookmark> => {
  try {
    const response = await api.put<Bookmark>(`/bookmarks/${id}`, data)
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to update bookmark')
    }
    throw error
  }
}

export const deleteBookmark = async (id: string): Promise<void> => {
  try {
    await api.delete(`/bookmarks/${id}`)
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to delete bookmark')
    }
    throw error
  }
}

export const searchBookmarks = async (
  query: string,
  page = 1,
  pageSize = 10
): Promise<PaginatedResponse<Bookmark>> => {
  try {
    const response = await api.get<PaginatedResponse<Bookmark>>('/bookmarks/search', {
      params: { query, page, pageSize },
    })
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to search bookmarks')
    }
    throw error
  }
}

export const getBookmarksByTag = async (
  tag: string,
  page = 1,
  pageSize = 10
): Promise<PaginatedResponse<Bookmark>> => {
  try {
    const response = await api.get<PaginatedResponse<Bookmark>>(`/bookmarks/tag/${tag}`, {
      params: { page, pageSize },
    })
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to get bookmarks by tag')
    }
    throw error
  }
}

export const getTags = async (): Promise<Record<string, number>> => {
  try {
    const response = await api.get<Record<string, number>>('/tags')
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Failed to get tags')
    }
    throw error
  }
}

export const healthCheck = async (): Promise<HealthCheckResponse> => {
  try {
    const response = await api.get<HealthCheckResponse>('/health')
    return response.data
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw new Error(error.response?.data?.message || 'Health check failed')
    }
    throw error
  }
}

export default api
