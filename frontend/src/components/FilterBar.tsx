'use client';

import React from 'react';
import { useBooks } from '../contexts/BookContext';

const FilterBar: React.FC = () => {
  const { 
    books, 
    selectedYear, 
    setSelectedYear, 
    sortBy, 
    setSortBy, 
    sortOrder, 
    setSortOrder,
    clearFilters,
    filteredBooks 
  } = useBooks();

  // Get unique years from books
  const availableYears = React.useMemo(() => {
    const years = books.map(book => book.year);
    const uniqueYears = Array.from(new Set(years)).sort((a, b) => b - a);
    return uniqueYears;
  }, [books]);

  const handleSortChange = (field: string) => {
    if (sortBy === field) {
      // Toggle sort order if same field
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      // Set new field with default desc order
      setSortBy(field);
      setSortOrder('desc');
    }
  };

  const getSortIcon = (field: string) => {
    if (sortBy !== field) {
      return (
        <svg className="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
        </svg>
      );
    }
    
    if (sortOrder === 'asc') {
      return (
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4l6 6m0 0l6-6m-6 6v12" />
        </svg>
      );
    } else {
      return (
        <svg className="w-4 h-4 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 20l6-6m0 0l6 6m-6-6V4" />
        </svg>
      );
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap gap-4 items-center">
        {/* Year Filter */}
        <div className="flex items-center space-x-2">
          <label htmlFor="year-filter" className="text-sm font-medium text-gray-700">
            Year:
          </label>
          <select
            id="year-filter"
            value={selectedYear}
            onChange={(e) => setSelectedYear(e.target.value)}
            className="border border-gray-300 rounded-md px-3 py-1 text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          >
            <option value="">All Years</option>
            {availableYears.map(year => (
              <option key={year} value={year.toString()}>
                {year}
              </option>
            ))}
          </select>
        </div>

        {/* Sort Options */}
        <div className="flex items-center space-x-2">
          <span className="text-sm font-medium text-gray-700">Sort by:</span>
          <div className="flex space-x-1">
            <button
              onClick={() => handleSortChange('title')}
              className={`flex items-center space-x-1 px-3 py-1 text-sm rounded-md transition-colors ${
                sortBy === 'title' 
                  ? 'bg-primary-100 text-primary-800' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              <span>Title</span>
              {getSortIcon('title')}
            </button>
            <button
              onClick={() => handleSortChange('author')}
              className={`flex items-center space-x-1 px-3 py-1 text-sm rounded-md transition-colors ${
                sortBy === 'author' 
                  ? 'bg-primary-100 text-primary-800' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              <span>Author</span>
              {getSortIcon('author')}
            </button>
            <button
              onClick={() => handleSortChange('year')}
              className={`flex items-center space-x-1 px-3 py-1 text-sm rounded-md transition-colors ${
                sortBy === 'year' 
                  ? 'bg-primary-100 text-primary-800' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              <span>Year</span>
              {getSortIcon('year')}
            </button>
            <button
              onClick={() => handleSortChange('created_at')}
              className={`flex items-center space-x-1 px-3 py-1 text-sm rounded-md transition-colors ${
                sortBy === 'created_at' 
                  ? 'bg-primary-100 text-primary-800' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              <span>Created</span>
              {getSortIcon('created_at')}
            </button>
          </div>
        </div>

        {/* Clear Filters */}
        <button
          onClick={clearFilters}
          className="text-sm text-primary-600 hover:text-primary-800 underline"
        >
          Clear Filters
        </button>
      </div>

      {/* Results Counter */}
      <div className="text-sm text-gray-600">
        Showing {filteredBooks.length} of {books.length} books
      </div>
    </div>
  );
};

export default FilterBar;