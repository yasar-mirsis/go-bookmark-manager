import React from 'react'
import type { Bookmark } from '../types'
import BookmarkCard from './BookmarkCard'
import styles from './BookmarkList.module.css'

interface BookmarkListProps {
  bookmarks: Bookmark[]
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  loading: boolean
  error?: string | null
  onBookmarkClick: (id: string) => void
  onEdit: (id: string) => void
  onDelete: (id: string) => void
  onTagClick: (tag: string) => void
}

function BookmarkList({
  bookmarks,
  currentPage,
  totalPages,
  onPageChange,
  loading,
  error,
  onBookmarkClick,
  onEdit,
  onDelete,
  onTagClick,
}: BookmarkListProps) {
  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1)
    }
  }

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1)
    }
  }

  const renderLoadingState = () => {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <p className={styles.loadingText}>Loading bookmarks...</p>
      </div>
    )
  }

  const renderErrorState = () => {
    return (
      <div className={styles.errorContainer}>
        <p className={styles.errorText}>Error: {error || 'An unexpected error occurred'}</p>
      </div>
    )
  }

  const renderEmptyState = () => {
    return (
      <div className={styles.emptyContainer}>
        <p className={styles.emptyText}>No bookmarks found</p>
      </div>
    )
  }

  const renderPagination = () => {
    if (totalPages <= 1) return null

    return (
      <div className={styles.pagination}>
        <button
          className={styles.paginationBtn}
          onClick={handlePrevious}
          disabled={currentPage <= 1}
          type="button"
        >
          Previous
        </button>
        <span className={styles.pageIndicator}>
          Page {currentPage} of {totalPages}
        </span>
        <button
          className={styles.paginationBtn}
          onClick={handleNext}
          disabled={currentPage >= totalPages}
          type="button"
        >
          Next
        </button>
      </div>
    )
  }

  const renderBookmarks = () => {
    return (
      <div className={styles.bookmarkGrid}>
        {bookmarks.map((bookmark) => (
          <BookmarkCard
            key={bookmark.id}
            bookmark={bookmark}
            onEdit={() => onEdit(bookmark.id)}
            onDelete={onDelete}
            onClick={onBookmarkClick}
            onTagClick={onTagClick}
          />
        ))}
      </div>
    )
  }

  if (loading) {
    return renderLoadingState()
  }

  if (error) {
    return renderErrorState()
  }

  if (bookmarks.length === 0) {
    return renderEmptyState()
  }

  return (
    <div className={styles.container}>
      {renderBookmarks()}
      {renderPagination()}
    </div>
  )
}

export default BookmarkList
