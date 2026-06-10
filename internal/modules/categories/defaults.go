package categories

func defaultCategories(userID uint) []Category {
	expense := CategoryTypeExpense
	income := CategoryTypeIncome

	return []Category{
		{UserID: userID, Name: "Food & Dining", Icon: "🍽️", Color: "#FF6B6B", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Transportation", Icon: "🚗", Color: "#4ECDC4", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Housing", Icon: "🏠", Color: "#45B7D1", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Healthcare", Icon: "🏥", Color: "#96CEB4", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Entertainment", Icon: "🎬", Color: "#FFEAA7", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Shopping", Icon: "🛍️", Color: "#DDA0DD", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Education", Icon: "📚", Color: "#98D8C8", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Utilities", Icon: "💡", Color: "#F7DC6F", Type: expense, IsDefault: true},
		{UserID: userID, Name: "Salary", Icon: "💰", Color: "#82E0AA", Type: income, IsDefault: true},
		{UserID: userID, Name: "Freelance", Icon: "💻", Color: "#85C1E9", Type: income, IsDefault: true},
		{UserID: userID, Name: "Other", Icon: "📦", Color: "#BDC3C7", Type: expense, IsDefault: true},
	}
}
