package options

//type IOption interface {
//	GetScope() func(db *gorm.DB) *gorm.DB
//	GetDependencies() map[string]func(db *gorm.DB) *gorm.DB
//}
//
//type Option struct {
//	Type         string
//	Dependencies map[string]func(db *gorm.DB) *gorm.DB
//}
//
//func (filter *Option) GetDependencies() map[string]func(db *gorm.DB) *gorm.DB {
//	return filter.Dependencies
//}
//
//func (filter *Option) GetScope() func(db *gorm.DB) *gorm.DB {
//	return filter.Scope
//}
