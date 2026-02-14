export interface PriceRecord {
  Ticker: string;
  Date: string;
  Price: number;
  Series: 'daily_price' | 'fair_value';
}
