import { render, screen } from '@testing-library/react';
import Header from './Header';
import '@testing-library/jest-dom';


describe('Header Component', () => {
  it('renders the application name', () => {
    render(Header());
    expect(screen.getByRole('heading', { name: /cibo/i })).toBeInTheDocument();
  });
});
