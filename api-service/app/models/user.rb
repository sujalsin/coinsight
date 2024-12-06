class User < ApplicationRecord
  has_secure_password
  
  # Associations
  has_one :portfolio, dependent: :destroy
  
  # Validations
  validates :email, presence: true, uniqueness: true,
            format: { with: URI::MailTo::EMAIL_REGEXP }
  validates :password, presence: true, length: { minimum: 6 }, if: :password_required?
  validates :username, presence: true, uniqueness: true,
            length: { minimum: 3, maximum: 30 }
  
  # Callbacks
  after_create :create_portfolio
  
  private
  
  def password_required?
    new_record? || password.present?
  end
  
  def create_portfolio
    Portfolio.create!(user: self)
  end
end
