import { useState, useEffect } from 'react'
import type { Bookmark, BookmarkFormData } from '../types'
import styles from './BookmarkForm.module.css'

interface BookmarkFormProps {
  initialValues?: Bookmark
  onSubmit: (data: BookmarkFormData) => void
  mode: 'create' | 'edit'
  onSubmitError?: (error: string) => void
}

interface FormErrors {
  url?: string
  title?: string
  description?: string
  tags?: string
}

function BookmarkForm({
  initialValues,
  onSubmit,
  mode,
  onSubmitError,
}: BookmarkFormProps) {
  const [formData, setFormData] = useState<BookmarkFormData>({
    url: '',
    title: '',
    description: '',
    tags: '',
  })
  const [errors, setErrors] = useState<FormErrors>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (initialValues) {
      setFormData({
        url: initialValues.url,
        title: initialValues.title,
        description: initialValues.description || '',
        tags: initialValues.tags.join(', '),
      })
    }
  }, [initialValues])

  const validateUrl = (url: string): string | undefined => {
    if (!url.trim()) {
      return 'URL is required'
    }
    try {
      new URL(url)
      return undefined
    } catch {
      return 'Please enter a valid URL (e.g., https://example.com)'
    }
  }

  const validateTitle = (title: string): string | undefined => {
    if (!title.trim()) {
      return 'Title is required'
    }
    if (title.trim().length < 2) {
      return 'Title must be at least 2 characters'
    }
    return undefined
  }

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {}

    const urlError = validateUrl(formData.url)
    if (urlError) {
      newErrors.url = urlError
    }

    const titleError = validateTitle(formData.title)
    if (titleError) {
      newErrors.title = titleError
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))

    // Clear error for the field being edited
    if (errors[name as keyof FormErrors]) {
      setErrors((prev) => ({
        ...prev,
        [name]: undefined,
      }))
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      return
    }

    setIsSubmitting(true)

    try {
      onSubmit(formData)
      // Reset form after successful submission
      if (mode === 'create') {
        setFormData({
          url: '',
          title: '',
          description: '',
          tags: '',
        })
      }
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : 'An unexpected error occurred'
      if (onSubmitError) {
        onSubmitError(errorMessage)
      } else {
        setErrors((prev) => ({
          ...prev,
          tags: errorMessage,
        }))
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.formGroup}>
        <label htmlFor="url" className={styles.label}>
          URL <span className={styles.required}>*</span>
        </label>
        <input
          type="text"
          id="url"
          name="url"
          value={formData.url}
          onChange={handleChange}
          placeholder="https://example.com"
          className={`${styles.input} ${errors.url ? styles.inputError : ''}`}
          disabled={isSubmitting}
          autoComplete="url"
        />
        {errors.url && <span className={styles.errorMessage}>{errors.url}</span>}
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="title" className={styles.label}>
          Title <span className={styles.required}>*</span>
        </label>
        <input
          type="text"
          id="title"
          name="title"
          value={formData.title}
          onChange={handleChange}
          placeholder="Enter a title for this bookmark"
          className={`${styles.input} ${errors.title ? styles.inputError : ''}`}
          disabled={isSubmitting}
          autoComplete="off"
        />
        {errors.title && (
          <span className={styles.errorMessage}>{errors.title}</span>
        )}
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="description" className={styles.label}>
          Description
        </label>
        <textarea
          id="description"
          name="description"
          value={formData.description}
          onChange={handleChange}
          placeholder="Add an optional description (max 500 characters)"
          className={styles.textarea}
          disabled={isSubmitting}
          maxLength={500}
          rows={3}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="tags" className={styles.label}>
          Tags
        </label>
        <input
          type="text"
          id="tags"
          name="tags"
          value={formData.tags}
          onChange={handleChange}
          placeholder="tag1, tag2, tag3 (comma-separated)"
          className={styles.input}
          disabled={isSubmitting}
          autoComplete="off"
        />
        <span className={styles.hint}>
          Enter tags separated by commas (e.g., programming, tutorials, news)
        </span>
      </div>

      {errors.tags && <span className={styles.errorMessage}>{errors.tags}</span>}

      <div className={styles.formActions}>
        <button
          type="submit"
          className={`${styles.submitButton} ${styles.btnPrimary}`}
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Saving...' : mode === 'create' ? 'Create Bookmark' : 'Update Bookmark'}
        </button>
      </div>
    </form>
  )
}

export default BookmarkForm
