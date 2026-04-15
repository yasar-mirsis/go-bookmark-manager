import React from 'react'
import type { Bookmark } from '../types'
import styles from './BookmarkCard.module.css'

interface BookmarkCardProps {
  bookmark: Bookmark
  onEdit: () => void
  onDelete: (id: string) => void
  onClick: (id: string) => void
  onTagClick: (tag: string) => void
}

function BookmarkCard({ bookmark, onEdit, onDelete, onClick, onTagClick }: BookmarkCardProps) {
  const truncateTitle = (title: string, maxLength: number = 50): string => {
    if (title.length <= maxLength) return title
    return title.substring(0, maxLength).trim() + '...'
  }

  const truncateDescription = (description: string, maxLength: number = 100): string => {
    if (!description) return ''
    if (description.length <= maxLength) return description
    return description.substring(0, maxLength).trim() + '...'
  }

  const truncateUrl = (url: string, maxLength: number = 50): string => {
    if (url.length <= maxLength) return url
    return url.substring(0, maxLength).trim() + '...'
  }

  const handleDelete = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (window.confirm('Are you sure you want to delete this bookmark?')) {
      onDelete(bookmark.id)
    }
  }

  const handleEdit = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    onEdit()
  }

  const handleTagClick = (e: React.MouseEvent, tag: string) => {
    e.preventDefault()
    e.stopPropagation()
    onTagClick(tag)
  }

  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <h3 className={styles.title} onClick={() => onClick(bookmark.id)} title={bookmark.title}>
          {truncateTitle(bookmark.title)}
        </h3>
      </div>

      <div className={styles.cardBody}>
        <p className={styles.url} title={bookmark.url}>
          {truncateUrl(bookmark.url)}
        </p>

        {bookmark.description && (
          <p className={styles.description} title={bookmark.description}>
            {truncateDescription(bookmark.description)}
          </p>
        )}

        {bookmark.tags.length > 0 && (
          <div className={styles.tags}>
            {bookmark.tags.map((tag) => (
              <button
                key={tag}
                className={styles.tagChip}
                onClick={(e) => handleTagClick(e, tag)}
                type="button"
              >
                {tag}
              </button>
            ))}
          </div>
        )}
      </div>

      <div className={styles.cardFooter}>
        <div className={styles.buttonGroup}>
          <button className={styles.btn} onClick={handleEdit} type="button">
            Edit
          </button>
          <button className={styles.btnDelete} onClick={handleDelete} type="button">
            Delete
          </button>
          <button className={styles.btnView} onClick={() => onClick(bookmark.id)} type="button">
            View
          </button>
        </div>
      </div>
    </div>
  )
}

export default BookmarkCard
