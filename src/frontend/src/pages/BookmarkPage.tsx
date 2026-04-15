import { useState, useEffect, useCallback } from 'react'
import api from '../services/api'
import type { Bookmark, TagInfo, BookmarkFormData, PaginatedResponse } from '../types'
import BookmarkList from '../components/BookmarkList'
import BookmarkModal from '../components/BookmarkModal'
import SearchBar from '../components/SearchBar'
import TagFilter from '../components/TagFilter'
import styles from './BookmarkPage.module.css'

function BookmarkPage() {
  const [bookmarks, setBookmarks] = useState<Bookmark[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedTag, setSelectedTag] = useState<string | undefined>(undefined)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [modalMode, setModalMode] = useState<'create' | 'edit'>('create')
  const [editingBookmark, setEditingBookmark] = useState<Bookmark | null>(null)
  const [tags, setTags] = useState<TagInfo[]>([])
  const pageSize = 10

  // Fetch bookmarks based on current filters
  const fetchBookmarks = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)

      let response: PaginatedResponse<Bookmark>

      if (searchQuery) {
        response = await api.searchBookmarks(searchQuery, currentPage, pageSize)
      } else if (selectedTag) {
        response = await api.getBookmarksByTag(selectedTag, currentPage, pageSize)
      } else {
        response = await api.getBookmarks(currentPage, pageSize)
      }

      setBookmarks(response.data)
      setTotalPages(response.totalPages)
    } catch (err) {
      setError('Failed to fetch bookmarks')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [currentPage, searchQuery, selectedTag])

  // Fetch tags
  const fetchTags = useCallback(async () => {
    try {
      const tagsData = await api.getTags()
      const tagInfo: TagInfo[] = Object.entries(tagsData).map(([name, count]) => ({
        name,
        count,
      }))
      setTags(tagInfo)
    } catch (err) {
      console.error('Failed to fetch tags:', err)
    }
  }, [])

  // Initial data fetch
  useEffect(() => {
    fetchBookmarks()
  }, [fetchBookmarks])

  useEffect(() => {
    fetchTags()
  }, [fetchTags])

  // Reset to page 1 when filters change
  useEffect(() => {
    setCurrentPage(1)
  }, [searchQuery, selectedTag])

  // Handle page change
  const handlePageChange = (page: number) => {
    setCurrentPage(page)
  }

  // Handle search
  const handleSearch = (query: string) => {
    setSearchQuery(query)
  }

  const handleClearSearch = () => {
    setSearchQuery('')
  }

  // Handle tag filter
  const handleTagClick = (tag: string) => {
    setSelectedTag(tag === selectedTag ? undefined : tag)
  }

  // Handle adding new bookmark
  const handleOpenCreateModal = () => {
    setModalMode('create')
    setEditingBookmark(null)
    setIsModalOpen(true)
  }

  // Handle editing bookmark
  const handleEditBookmark = (id: string) => {
    const bookmark = bookmarks.find((b) => b.id === id)
    if (bookmark) {
      setModalMode('edit')
      setEditingBookmark(bookmark)
      setIsModalOpen(true)
    }
  }

  // Handle deleting bookmark
  const handleDeleteBookmark = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this bookmark?')) {
      return
    }

    try {
      await api.deleteBookmark(id)
      setBookmarks((prev) => prev.filter((b) => b.id !== id))

      // If we're on the last page and it's now empty, go back one page
      if (bookmarks.length === 1 && currentPage > 1) {
        setCurrentPage(currentPage - 1)
      }
    } catch (err) {
      setError('Failed to delete bookmark')
      console.error(err)
    }
  }

  // Handle bookmark click (navigate to detail page)
  const handleBookmarkClick = (id: string) => {
    window.location.href = `/bookmarks/${id}`
  }

  // Handle form submission for create
  const handleCreateBookmark = async (data: BookmarkFormData) => {
    try {
      const tagsArray = data.tags
        .split(',')
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0)

      await api.createBookmark({
        url: data.url,
        title: data.title,
        description: data.description,
        tags: tagsArray,
      })

      // Refresh bookmarks list
      await fetchBookmarks()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create bookmark'
      throw errorMessage
    }
  }

  // Handle form submission for edit
  const handleUpdateBookmark = async (data: BookmarkFormData) => {
    if (!editingBookmark) return

    try {
      const tagsArray = data.tags
        .split(',')
        .map((tag) => tag.trim())
        .filter((tag) => tag.length > 0)

      await api.updateBookmark(editingBookmark.id, {
        url: data.url,
        title: data.title,
        description: data.description,
        tags: tagsArray,
      })

      // Refresh bookmarks list
      await fetchBookmarks()
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to update bookmark'
      throw errorMessage
    }
  }

  // Handle form submission based on mode
  const handleFormSubmit = async (data: BookmarkFormData) => {
    if (modalMode === 'create') {
      await handleCreateBookmark(data)
    } else {
      await handleUpdateBookmark(data)
    }
  }

  // Handle form submission error
  const handleFormSubmitError = (errorMessage: string) => {
    setError(errorMessage)
  }

  // Close modal
  const handleCloseModal = () => {
    setIsModalOpen(false)
    setEditingBookmark(null)
  }

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <h1 className={styles.title}>Bookmarks</h1>
        <button
          className={styles.addButton}
          onClick={handleOpenCreateModal}
          type="button"
        >
          Add Bookmark
        </button>
      </header>

      <div className={styles.searchContainer}>
        <SearchBar onSearch={handleSearch} onClear={handleClearSearch} />
      </div>

      <div className={styles.content}>
        <aside className={styles.sidebar}>
          <TagFilter
            tags={tags}
            selectedTag={selectedTag}
            onTagClick={handleTagClick}
          />
        </aside>

        <main className={styles.main}>
          {error && (
            <div className={styles.errorBanner}>
              {error}
              <button
                className={styles.dismissButton}
                onClick={() => setError(null)}
                type="button"
              >
                ×
              </button>
            </div>
          )}

          <BookmarkList
            bookmarks={bookmarks}
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
            loading={loading}
            error={error}
            onBookmarkClick={handleBookmarkClick}
            onEdit={handleEditBookmark}
            onDelete={handleDeleteBookmark}
            onTagClick={handleTagClick}
          />
        </main>
      </div>

      <BookmarkModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={handleFormSubmit}
        bookmark={editingBookmark}
        mode={modalMode}
        onSubmitError={handleFormSubmitError}
      />
    </div>
  )
}

export default BookmarkPage
