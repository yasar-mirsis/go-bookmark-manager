import { useEffect, useRef } from 'react'
import BookmarkForm from './BookmarkForm'
import type { Bookmark, BookmarkFormData } from '../types'
import styles from './BookmarkModal.module.css'

interface BookmarkModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (data: BookmarkFormData) => void
  bookmark?: Bookmark | null
  mode: 'create' | 'edit'
  onSubmitError?: (error: string) => void
}

function BookmarkModal({
  isOpen,
  onClose,
  onSubmit,
  bookmark,
  mode,
  onSubmitError,
}: BookmarkModalProps) {
  const modalRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose()
      }
    }

    document.addEventListener('keydown', handleEscape)
    return () => {
      document.removeEventListener('keydown', handleEscape)
    }
  }, [isOpen, onClose])

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }

    return () => {
      document.body.style.overflow = ''
    }
  }, [isOpen])

  const handleOverlayClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (modalRef.current && !modalRef.current.contains(e.target as Node)) {
      onClose()
    }
  }

  const handleFormSubmit = (data: BookmarkFormData) => {
    onSubmit(data)
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className={styles.overlay} onClick={handleOverlayClick} role="dialog" aria-modal="true">
      <div className={styles.modal} ref={modalRef}>
        <div className={styles.header}>
          <h2 className={styles.title}>
            {mode === 'create' ? 'Add Bookmark' : 'Edit Bookmark'}
          </h2>
          <button
            className={styles.closeButton}
            onClick={onClose}
            aria-label="Close modal"
            type="button"
          >
            ×
          </button>
        </div>
        <div className={styles.content}>
          <BookmarkForm
            initialValues={bookmark || undefined}
            onSubmit={handleFormSubmit}
            mode={mode}
            onSubmitError={onSubmitError}
          />
        </div>
      </div>
    </div>
  )
}

export default BookmarkModal
