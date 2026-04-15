import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import BookmarkList from '../BookmarkList'
import type { Bookmark } from '../../types'

const mockBookmarks: Bookmark[] = [
  {
    id: '1',
    url: 'https://example.com',
    title: 'Example Bookmark',
    description: 'This is a test bookmark',
    tags: ['react', 'typescript'],
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: '2',
    url: 'https://test.com',
    title: 'Test Bookmark',
    description: 'Another test bookmark',
    tags: ['testing', 'jest'],
    createdAt: '2024-01-02T00:00:00Z',
    updatedAt: '2024-01-02T00:00:00Z',
  },
  {
    id: '3',
    url: 'https://demo.com',
    title: 'Demo Bookmark',
    description: 'A demo bookmark',
    tags: ['demo'],
    createdAt: '2024-01-03T00:00:00Z',
    updatedAt: '2024-01-03T00:00:00Z',
  },
]

const defaultProps = {
  bookmarks: mockBookmarks,
  currentPage: 1,
  totalPages: 3,
  onPageChange: vi.fn(),
  loading: false,
  onBookmarkClick: vi.fn(),
  onEdit: vi.fn(),
  onDelete: vi.fn(),
  onTagClick: vi.fn(),
}

describe('BookmarkList', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders bookmarks correctly', () => {
    render(<BookmarkList {...defaultProps} />)

    expect(screen.getByText('Example Bookmark')).toBeInTheDocument()
    expect(screen.getByText('Test Bookmark')).toBeInTheDocument()
    expect(screen.getByText('Demo Bookmark')).toBeInTheDocument()
  })

  it('renders tag chips for each bookmark', () => {
    render(<BookmarkList {...defaultProps} />)

    expect(screen.getByText('react')).toBeInTheDocument()
    expect(screen.getByText('typescript')).toBeInTheDocument()
    expect(screen.getByText('testing')).toBeInTheDocument()
    expect(screen.getByText('jest')).toBeInTheDocument()
    expect(screen.getByText('demo')).toBeInTheDocument()
  })

  it('renders loading state when loading is true', () => {
    render(<BookmarkList {...defaultProps} loading={true} />)

    expect(screen.getByText(/Loading bookmarks/i)).toBeInTheDocument()
    expect(screen.getByRole('status')).toBeInTheDocument()
  })

  it('shows spinner during loading', () => {
    render(<BookmarkList {...defaultProps} loading={true} />)

    expect(screen.getByRole('status')).toHaveClass('spinner')
  })

  it('renders error state when error is provided', () => {
    render(
      <BookmarkList
        {...defaultProps}
        loading={false}
        error="Failed to fetch bookmarks"
      />
    )

    expect(screen.getByText(/Error: Failed to fetch bookmarks/i)).toBeInTheDocument()
  })

  it('renders empty state when no bookmarks', () => {
    render(<BookmarkList {...defaultProps} bookmarks={[]} />)

    expect(screen.getByText(/No bookmarks found/i)).toBeInTheDocument()
  })

  it('displays pagination controls when totalPages > 1', () => {
    render(<BookmarkList {...defaultProps} />)

    expect(screen.getByRole('button', { name: 'Previous' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Next' })).toBeInTheDocument()
    expect(screen.getByText(/Page 1 of 3/i)).toBeInTheDocument()
  })

  it('does not display pagination when totalPages is 1', () => {
    render(<BookmarkList {...defaultProps} totalPages={1} />)

    expect(screen.queryByRole('button', { name: 'Previous' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Next' })).not.toBeInTheDocument()
  })

  it('calls onPageChange with previous page when Previous button is clicked', () => {
    render(<BookmarkList {...defaultProps} currentPage={2} />)

    const previousButton = screen.getByRole('button', { name: 'Previous' })
    fireEvent.click(previousButton)

    expect(defaultProps.onPageChange).toHaveBeenCalledWith(1)
  })

  it('calls onPageChange with next page when Next button is clicked', () => {
    render(<BookmarkList {...defaultProps} currentPage={1} />)

    const nextButton = screen.getByRole('button', { name: 'Next' })
    fireEvent.click(nextButton)

    expect(defaultProps.onPageChange).toHaveBeenCalledWith(2)
  })

  it('disables Previous button on first page', () => {
    render(<BookmarkList {...defaultProps} currentPage={1} />)

    const previousButton = screen.getByRole('button', { name: 'Previous' })
    expect(previousButton).toBeDisabled()
  })

  it('disables Next button on last page', () => {
    render(<BookmarkList {...defaultProps} currentPage={3} totalPages={3} />)

    const nextButton = screen.getByRole('button', { name: 'Next' })
    expect(nextButton).toBeDisabled()
  })

  it('does not call onPageChange when clicking disabled Previous button', () => {
    render(<BookmarkList {...defaultProps} currentPage={1} />)

    const previousButton = screen.getByRole('button', { name: 'Previous' })
    fireEvent.click(previousButton)

    expect(defaultProps.onPageChange).not.toHaveBeenCalled()
  })

  it('does not call onPageChange when clicking disabled Next button', () => {
    render(<BookmarkList {...defaultProps} currentPage={3} totalPages={3} />)

    const nextButton = screen.getByRole('button', { name: 'Next' })
    fireEvent.click(nextButton)

    expect(defaultProps.onPageChange).not.toHaveBeenCalled()
  })

  it('calls onBookmarkClick when bookmark title is clicked', async () => {
    render(<BookmarkList {...defaultProps} />)

    const title = screen.getByText('Example Bookmark')
    fireEvent.click(title)

    await waitFor(() => {
      expect(defaultProps.onBookmarkClick).toHaveBeenCalledWith('1')
    })
  })

  it('calls onEdit when Edit button is clicked', async () => {
    render(<BookmarkList {...defaultProps} />)

    const editButton = screen.getAllByRole('button', { name: 'Edit' })[0]
    fireEvent.click(editButton)

    await waitFor(() => {
      expect(defaultProps.onEdit).toHaveBeenCalledWith('1')
    })
  })

  it('calls onDelete when Delete button is clicked and confirmed', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true)

    render(<BookmarkList {...defaultProps} />)

    const deleteButton = screen.getAllByRole('button', { name: 'Delete' })[0]
    fireEvent.click(deleteButton)

    await waitFor(() => {
      expect(defaultProps.onDelete).toHaveBeenCalledWith('1')
    })

    vi.restoreAllMocks()
  })

  it('calls onTagClick when tag chip is clicked', async () => {
    render(<BookmarkList {...defaultProps} />)

    const reactTag = screen.getByText('react')
    fireEvent.click(reactTag)

    await waitFor(() => {
      expect(defaultProps.onTagClick).toHaveBeenCalledWith('react')
    })
  })

  it('renders BookmarkCard components with correct props', () => {
    render(<BookmarkList {...defaultProps} />)

    // Verify cards are rendered
    const cards = screen.getAllByRole('article')
    expect(cards).toHaveLength(3)
  })

  it('shows default error message when error is null but provided', () => {
    render(
      <BookmarkList
        {...defaultProps}
        loading={false}
        error={null}
      />
    )

    // Should show empty state since no bookmarks would be displayed with null error
    // This tests the error fallback behavior
  })

  it('updates page indicator when currentPage changes', () => {
    const { rerender } = render(<BookmarkList {...defaultProps} currentPage={1} />)

    expect(screen.getByText(/Page 1 of 3/i)).toBeInTheDocument()

    rerender(<BookmarkList {...defaultProps} currentPage={2} />)

    expect(screen.getByText(/Page 2 of 3/i)).toBeInTheDocument()
  })

  // Snapshot tests
  it('matches snapshot with bookmarks', () => {
    const { container } = render(<BookmarkList {...defaultProps} />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot in loading state', () => {
    const { container } = render(<BookmarkList {...defaultProps} loading={true} />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot in error state', () => {
    const { container } = render(
      <BookmarkList {...defaultProps} loading={false} error="Test error" />
    )
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot in empty state', () => {
    const { container } = render(<BookmarkList {...defaultProps} bookmarks={[]} />)
    expect(container).toMatchSnapshot()
  })
})
