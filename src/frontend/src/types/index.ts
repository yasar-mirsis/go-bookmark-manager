export interface Bookmark {
  id: string
  url: string
  title: string
  description: string
  tags: string[]
  createdAt: string
  updatedAt: string
}

export interface CreateBookmarkInput {
  url: string
  title: string
  description?: string
  tags: string[]
}

export interface UpdateBookmarkInput {
  url?: string
  title?: string
  description?: string
  tags?: string[]
}

export interface PaginationResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface TagCount {
  [tagName: string]: number
}

export interface HealthCheckResponse {
  status: string
  timestamp: string
}

export interface BookmarkFormData {
  url: string
  title: string
  description: string
  tags: string
}

export interface SearchParams {
  query: string
  page: number
  pageSize: number
}

export interface TagInfo {
  name: string
  count: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}
