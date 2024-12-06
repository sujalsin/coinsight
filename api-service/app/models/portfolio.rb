class Portfolio < ApplicationRecord
  belongs_to :user
  has_many :holdings, dependent: :destroy
  
  # Store historical returns as a JSON array
  serialize :returns, Array
  
  # Validations
  validates :user_id, presence: true, uniqueness: true
  
  def total_value
    holdings.sum { |holding| holding.amount * holding.current_price }
  end
  
  def risk_level
    return 'low' if returns.empty?
    
    volatility = calculate_volatility
    case volatility
    when 0..0.1 then 'low'
    when 0.1..0.2 then 'medium'
    else 'high'
    end
  end
  
  private
  
  def calculate_volatility
    return 0 if returns.size < 2
    
    mean = returns.sum / returns.size.to_f
    variance = returns.sum { |r| (r - mean) ** 2 } / (returns.size - 1)
    Math.sqrt(variance)
  end
end
