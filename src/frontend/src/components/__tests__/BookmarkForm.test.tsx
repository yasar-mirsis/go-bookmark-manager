import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import BookmarkForm from '../BookmarkForm'
import type { Bookmark, BookmarkFormData } from '../../types'

const mockInitialValues: Bookmark = {
  id: '1',
  url: 'https://example.com',
  title: 'Example Bookmark',
  description: 'This is a test description',
  tags: ['react', 'typescript'],
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
}

const defaultProps = {
  onSubmit: vi.fn(),
  mode: 'create' as const,
}

describe('BookmarkForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders form fields correctly', () => {
    render(<BookmarkForm {...defaultProps} />)

    expect(screen.getByLabelText(/URL \*/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Title \*/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Description/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/Tags/i)).toBeInTheDocument()
  })

  it('displays placeholder text for all fields', () => {
    render(<BookmarkForm {...defaultProps} />)

    expect(screen.getByPlaceholderText(/https:\/\/example.com/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/Enter a title/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/Add an optional description/i)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(/tag1, tag2, tag3/i)).toBeInTheDocument()
  })

  it('fills form with initial values in edit mode', () => {
    render(
      <BookmarkForm
        initialValues={mockInitialValues}
        onSubmit={defaultProps.onSubmit}
        mode="edit"
      />
    )

    expect(screen.getByLabelText(/URL \*/i)).toHaveValue('https://example.com')
    expect(screen.getByLabelText(/Title \*/i)).toHaveValue('Example Bookmark')
    expect(screen.getByLabelText(/Description/i)).toHaveValue('This is a test description')
    expect(screen.getByLabelText(/Tags/i)).toHaveValue('react, typescript')
  })

  it('validates URL is required', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test Title' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/URL is required/i)).toBeInTheDocument()
    })
  })

  it('validates URL format', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'invalid-url' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test Title' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/Please enter a valid URL/i)).toBeInTheDocument()
    })
  })

  it('validates title is required', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/Title is required/i)).toBeInTheDocument()
    })
  })

  it('validates title minimum length', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'A' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/Title must be at least 2 characters/i)).toBeInTheDocument()
    })
  })

  it('clears validation errors when user starts typing', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(screen.getByText(/URL is required/i)).toBeInTheDocument()
    })

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    expect(screen.queryByText(/URL is required/i)).not.toBeInTheDocument()
  })

  it('submits form with valid data', async () => {
    const mockSubmit = vi.fn()
    render(<BookmarkForm onSubmit={mockSubmit} mode="create" />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test Bookmark' } })

    const descInput = screen.getByLabelText(/Description/i)
    fireEvent.change(descInput, { target: { value: 'Test description' } })

    const tagsInput = screen.getByLabelText(/Tags/i)
    fireEvent.change(tagsInput, { target: { value: 'tag1, tag2' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalledWith({
        url: 'https://example.com',
        title: 'Test Bookmark',
        description: 'Test description',
        tags: 'tag1, tag2',
      })
    })
  })

  it('submits form without description (optional field)', async () => {
    const mockSubmit = vi.fn()
    render(<BookmarkForm onSubmit={mockSubmit} mode="create" />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test Bookmark' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalledWith({
        url: 'https://example.com',
        title: 'Test Bookmark',
        description: '',
        tags: '',
      })
    })
  })

  it('displays error message when onSubmit throws', async () => {
    const mockSubmit = vi.fn().mockImplementation(() => {
      throw new Error('Network error')
    })
    const mockErrorCallback = vi.fn()

    render(
      <BookmarkForm
        onSubmit={mockSubmit}
        mode="create"
        onSubmitError={mockErrorCallback}
      />
    )

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(mockErrorCallback).toHaveBeenCalledWith('Network error')
    })
  })

  it('shows disabled state during submit', async () => {
    let resolvePromise: () => void
    const mockSubmit = vi.fn(
      () =>
        new Promise<void>((resolve) => {
          resolvePromise = resolve
        })
    )

    render(<BookmarkForm onSubmit={mockSubmit} mode="create" />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test' } })

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })

    fireEvent.click(submitButton)

    await waitFor(() => {
      expect(submitButton).toBeDisabled()
      expect(submitButton).toHaveTextContent('Saving...')
    })

    // Resolve the promise
    resolvePromise!()
  })

  it('disables all inputs during submit', async () => {
    let resolvePromise: () => void
    const mockSubmit = vi.fn(
      () =>
        new Promise<void>((resolve) => {
          resolvePromise = resolve
        })
    )

    render(<BookmarkForm onSubmit={mockSubmit} mode="create" />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    const titleInput = screen.getByLabelText(/Title \*/i)
    const descInput = screen.getByLabelText(/Description/i)
    const tagsInput = screen.getByLabelText(/Tags/i)

    fireEvent.click(screen.getByRole('button', { name: /Create Bookmark/i }))

    await waitFor(() => {
      expect(urlInput).toBeDisabled()
      expect(titleInput).toBeDisabled()
      expect(descInput).toBeDisabled()
      expect(tagsInput).toBeDisabled()
    })

    resolvePromise!()
  })

  it('resets form after successful submission in create mode', async () => {
    const mockSubmit = vi.fn()
    render(<BookmarkForm onSubmit={mockSubmit} mode="create" />)

    const urlInput = screen.getByLabelText(/URL \*/i)
    fireEvent.change(urlInput, { target: { value: 'https://example.com' } })

    const titleInput = screen.getByLabelText(/Title \*/i)
    fireEvent.change(titleInput, { target: { value: 'Test' } })

    fireEvent.click(screen.getByRole('button', { name: /Create Bookmark/i }))

    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalled()
    })

    expect(urlInput).toHaveValue('')
    expect(titleInput).toHaveValue('')
  })

  it('does not reset form after submission in edit mode', async () => {
    const mockSubmit = vi.fn()
    render(
      <BookmarkForm
        initialValues={mockInitialValues}
        onSubmit={mockSubmit}
        mode="edit"
      />
    )

    fireEvent.click(screen.getByRole('button', { name: /Update Bookmark/i }))

    await waitFor(() => {
      expect(mockSubmit).toHaveBeenCalled()
    })

    expect(screen.getByLabelText(/URL \*/i)).toHaveValue('https://example.com')
    expect(screen.getByLabelText(/Title \*/i)).toHaveValue('Example Bookmark')
  })

  it('shows hint text for tags field', () => {
    render(<BookmarkForm {...defaultProps} />)
    expect(screen.getByText(/Enter tags separated by commas/i)).toBeInTheDocument()
  })

  it('applies error styling to input fields with errors', async () => {
    render(<BookmarkForm {...defaultProps} />)

    const submitButton = screen.getByRole('button', { name: /Create Bookmark/i })
    fireEvent.click(submitButton)

    await waitFor(() => {
      const urlInput = screen.getByLabelText(/URL \*/i)
      expect(urlInput).toHaveClass('inputError')
    })
  })

  // Snapshot test
  it('matches snapshot in create mode', () => {
    const { container } = render(<BookmarkForm {...defaultProps} />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot in edit mode', () => {
    const { container } = render(
      <BookmarkForm
        initialValues={mockInitialValues}
        onSubmit={defaultProps.onSubmit}
        mode="edit"
      />
    )
    expect(container).toMatchSnapshot()
  })
})
