import { render, screen, fireEvent } from '@testing-library/react';
import Sidebar from './Sidebar';

describe('Sidebar Component', () => {
  it('renders placeholder text when open', () => {
    const mockToggle = jest.fn();
    render(Sidebar({ isOpen: true, onToggle: mockToggle }));
    expect(screen.getByText(/data sources/i)).toBeInTheDocument();
  });

  it('hides content when collapsed', () => {
    const mockToggle = jest.fn();
    render(Sidebar({ isOpen: false, onToggle: mockToggle }));
    expect(screen.queryByText(/data sources/i)).not.toBeInTheDocument();
  });

  it('calls onToggle when the toggle button is clicked', () => {
    const mockToggle = jest.fn();
    render(Sidebar({ isOpen: true, onToggle: mockToggle }));
    fireEvent.click(screen.getByRole('button', { name: /toggle sidebar/i }));
    expect(mockToggle).toHaveBeenCalledTimes(1);
  });
});
