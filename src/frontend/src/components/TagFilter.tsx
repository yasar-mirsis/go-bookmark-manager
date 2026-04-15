import styles from './TagFilter.module.css'
import type { TagInfo } from '../types'

interface TagFilterProps {
  tags: TagInfo[]
  selectedTag?: string
  onTagClick: (tag: string) => void
}

function TagFilter({ tags, selectedTag, onTagClick }: TagFilterProps) {
  const handleClearFilter = () => {
    onTagClick('')
  }

  if (tags.length === 0) {
    return null
  }

  return (
    <div className={styles.tagFilter}>
      <h3 className={styles.title}>Filter by Tag</h3>
      <div className={styles.tagList}>
        {tags.map((tag) => (
          <button
            key={tag.name}
            onClick={() => onTagClick(tag.name)}
            className={`${styles.tagChip} ${selectedTag === tag.name ? styles.selected : ''}`}
            type="button"
            aria-pressed={selectedTag === tag.name}
          >
            <span className={styles.tagName}>{tag.name}</span>
            <span className={styles.tagCount}>({tag.count})</span>
          </button>
        ))}
      </div>
      {selectedTag && (
        <button
          onClick={handleClearFilter}
          className={styles.clearButton}
          type="button"
        >
          Clear Filter
        </button>
      )}
    </div>
  )
}

export default TagFilter
