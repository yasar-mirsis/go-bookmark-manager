import { render, screen, fireEvent } from '@testing-library/react'
import '@testing-library/jest-dom'
import TagFilter from '../TagFilter'
import type { TagInfo } from '../../types'

const mockTags: TagInfo[] = [
  { name: 'react', count: 15 },
  { name: 'typescript', count: 8 },
  { name: 'javascript', count: 23 },
  { name: 'testing', count: 5 },
  { name: 'devops', count: 3 },
]

const defaultProps = {
  tags: mockTags,
  onTagClick: vi.fn(),
}

describe('TagFilter', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders title correctly', () => {
    render(<TagFilter {...defaultProps} />)

    expect(screen.getByText('Filter by Tag')).toBeInTheDocument()
  })

  it('renders all tags with their counts', () => {
    render(<TagFilter {...defaultProps} />)

    expect(screen.getByText('react')).toBeInTheDocument()
    expect(screen.getByText('(15)')).toBeInTheDocument()
    expect(screen.getByText('typescript')).toBeInTheDocument()
    expect(screen.getByText('(8)')).toBeInTheDocument()
    expect(screen.getByText('javascript')).toBeInTheDocument()
    expect(screen.getByText('(23)')).toBeInTheDocument()
    expect(screen.getByText('testing')).toBeInTheDocument()
    expect(screen.getByText('(5)')).toBeInTheDocument()
    expect(screen.getByText('devops')).toBeInTheDocument()
    expect(screen.getByText('(3)')).toBeInTheDocument()
  })

  it('renders each tag as a button', () => {
    render(<TagFilter {...defaultProps} />)

    const reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    expect(reactButton).toBeInTheDocument()

    const typescriptButton = screen.getByRole('button', { name: /typescript \(8\)/i })
    expect(typescriptButton).toBeInTheDocument()
  })

  it('calls onTagClick when a tag is clicked', () => {
    render(<TagFilter {...defaultProps} />)

    const reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    fireEvent.click(reactButton)

    expect(defaultProps.onTagClick).toHaveBeenCalledWith('react')
  })

  it('calls onTagClick with correct tag name for each tag', () => {
    render(<TagFilter {...defaultProps} />)

    const tags = ['react', 'typescript', 'javascript', 'testing', 'devops']

    tags.forEach((tag) => {
      const button = screen.getByRole('button', { name: new RegExp(tag, 'i') })
      fireEvent.click(button)
      expect(defaultProps.onTagClick).toHaveBeenCalledWith(tag)
    })
  })

  it('applies selected class to selected tag', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    expect(reactButton).toHaveClass('selected')
  })

  it('does not apply selected class to non-selected tags', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const typescriptButton = screen.getByRole('button', { name: /typescript \(8\)/i })
    expect(typescriptButton).not.toHaveClass('selected')
  })

  it('updates selected tag when selection changes', () => {
    const { rerender } = render(<TagFilter {...defaultProps} selectedTag="react" />)

    let reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    expect(reactButton).toHaveClass('selected')

    rerender(<TagFilter {...defaultProps} selectedTag="typescript" />)

    reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    expect(reactButton).not.toHaveClass('selected')

    const typescriptButton = screen.getByRole('button', { name: /typescript \(8\)/i })
    expect(typescriptButton).toHaveClass('selected')
  })

  it('renders clear filter button when a tag is selected', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    expect(screen.getByRole('button', { name: 'Clear Filter' })).toBeInTheDocument()
  })

  it('does not render clear filter button when no tag is selected', () => {
    render(<TagFilter {...defaultProps} />)

    expect(screen.queryByRole('button', { name: 'Clear Filter' })).not.toBeInTheDocument()
  })

  it('calls onTagClick with empty string when clear filter is clicked', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const clearButton = screen.getByRole('button', { name: 'Clear Filter' })
    fireEvent.click(clearButton)

    expect(defaultProps.onTagClick).toHaveBeenCalledWith('')
  })

  it('sets aria-pressed to true for selected tag', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    expect(reactButton).toHaveAttribute('aria-pressed', 'true')
  })

  it('sets aria-pressed to false for non-selected tags', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const typescriptButton = screen.getByRole('button', { name: /typescript \(8\)/i })
    expect(typescriptButton).toHaveAttribute('aria-pressed', 'false')
  })

  it('does not render anything when tags array is empty', () => {
    const { container } = render(<TagFilter {...defaultProps} tags={[]} />)

    expect(container.firstChild).toBeNull()
  })

  it('handles tags with special characters in name', () => {
    const specialTags: TagInfo[] = [
      { name: 'react-native', count: 5 },
      { name: 'node.js', count: 10 },
      { name: 'c#', count: 3 },
    ]

    render(<TagFilter {...defaultProps} tags={specialTags} />)

    expect(screen.getByText('react-native')).toBeInTheDocument()
    expect(screen.getByText('node.js')).toBeInTheDocument()
    expect(screen.getByText('c#')).toBeInTheDocument()
  })

  it('handles tags with zero count', () => {
    const tagsWithZero: TagInfo[] = [
      { name: 'empty-tag', count: 0 },
    ]

    render(<TagFilter {...defaultProps} tags={tagsWithZero} />)

    expect(screen.getByText('empty-tag')).toBeInTheDocument()
    expect(screen.getByText('(0)')).toBeInTheDocument()
  })

  it('handles large tag counts correctly', () => {
    const tagsWithLargeCount: TagInfo[] = [
      { name: 'popular', count: 1234 },
    ]

    render(<TagFilter {...defaultProps} tags={tagsWithLargeCount} />)

    expect(screen.getByText('popular')).toBeInTheDocument()
    expect(screen.getByText('(1234)')).toBeInTheDocument()
  })

  it('tags are rendered in the correct order', () => {
    render(<TagFilter {...defaultProps} />)

    const tagButtons = screen.getAllByRole('button')
    // First 5 buttons should be the tags, last one is clear filter (if selected)
    expect(tagButtons[0]).toHaveTextContent('react')
    expect(tagButtons[1]).toHaveTextContent('typescript')
    expect(tagButtons[2]).toHaveTextContent('javascript')
    expect(tagButtons[3]).toHaveTextContent('testing')
    expect(tagButtons[4]).toHaveTextContent('devops')
  })

  it('clicking same selected tag still calls onTagClick', () => {
    render(<TagFilter {...defaultProps} selectedTag="react" />)

    const reactButton = screen.getByRole('button', { name: /react \(15\)/i })
    fireEvent.click(reactButton)

    expect(defaultProps.onTagClick).toHaveBeenCalledWith('react')
  })

  // Snapshot tests
  it('matches snapshot with no selection', () => {
    const { container } = render(<TagFilter {...defaultProps} />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot with selection', () => {
    const { container } = render(<TagFilter {...defaultProps} selectedTag="react" />)
    expect(container).toMatchSnapshot()
  })

  it('matches snapshot with single tag', () => {
    const { container } = render(
      <TagFilter
        {...defaultProps}
        tags={[{ name: 'single', count: 1 }]}
      />
    )
    expect(container).toMatchSnapshot()
  })
})
