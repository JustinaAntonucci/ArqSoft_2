package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	controller "hotels-api/controllers/hotel"
	dto "hotels-api/dtos"
	service "hotels-api/services"
	"hotels-api/utils/errors"
)

// MockHotelService es el mock del servicio de hotel
type MockHotelService struct {
	mock.Mock
}

func (m *MockHotelService) GetHotel(id string) (dto.HotelDto, errors.ApiError) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return dto.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dto.HotelDto), nil
}

func (m *MockHotelService) InsertHotel(hotel dto.HotelDto) (dto.HotelDto, errors.ApiError) {
	args := m.Called(hotel)
	if args.Get(0) == nil {
		return dto.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dto.HotelDto), nil
}

func (m *MockHotelService) UpdateHotel(id string, hotel dto.HotelDto) (dto.HotelDto, errors.ApiError) {
	args := m.Called(id, hotel)
	if args.Get(0) == nil {
		return dto.HotelDto{}, args.Get(1).(errors.ApiError)
	}
	return args.Get(0).(dto.HotelDto), nil
}
func (m *MockHotelService) DeleteHotel(id string) errors.ApiError {
	args := m.Called(id)
	errArg := args.Get(0)
	if errArg != nil {
		if apiErr, ok := errArg.(errors.ApiError); ok {
			return apiErr
		} else {
			// Aquí puedes manejar la conversión de tipos si es necesario
			return errors.NewInternalServerApiError("Error casting to ApiError", nil)
		}
	}
	return nil
}

func TestGet(t *testing.T) {
	mockService := new(MockHotelService)
	service.HotelService = mockService

	expectedHotel := dto.HotelDto{
		Id:          "example-id",
		Name:        "MockHotel",
		Description: "Mock Description",
		Country:     "Mock Country",
		City:        "Mock City",
		Adress:      "Mock Address",
		Images:      []string{"image1.jpg", "image2.jpg"},
		Amenities:   []string{"wifi", "pool"},
	}

	mockService.On("GetHotel", "example-id").Return(expectedHotel, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "example-id"}}

	controller.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response dto.HotelDto
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedHotel, response)
	mockService.AssertExpectations(t)
}

func TestInsert(t *testing.T) {
	// Mock del servicio de hotel
	mockService := new(MockHotelService)
	service.HotelService = mockService

	// Datos de ejemplo del hotel recibidos del contexto
	hotelDto := dto.HotelDto{
		Name:        "New Hotel",
		Description: "New Description",
		Country:     "New Country",
		City:        "New City",
		Adress:      "New Address",
		Images:      []string{"new_image1.jpg", "new_image2.jpg"},
		Amenities:   []string{"wifi", "pool"},
	}

	// Mock del servicio esperando recibir el nuevo hotel DTO
	mockService.On("InsertHotel", hotelDto).Return(hotelDto, nil)

	// Simulación de la solicitud HTTP con Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/hotels", strings.NewReader(`{
        "Name": "New Hotel",
        "Description": "New Description",
        "Country": "New Country",
        "City": "New City",
        "Address": "New Address",
        "Images": ["new_image1.jpg", "new_image2.jpg"],
        "Amenities": ["wifi", "pool"]
    }`))
	c.Request.Header.Set("Content-Type", "application/json")

	// Llamada al controlador con el contexto creado
	controller.Insert(c)

	// Verificar que la respuesta sea la esperada
	assert.Equal(t, http.StatusCreated, w.Code)
	var response dto.HotelDto
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON response: %v", err)
	}

	// Verificar que los datos devueltos coincidan con los datos esperados
	assert.Equal(t, hotelDto.Name, response.Name)
	assert.Equal(t, hotelDto.Description, response.Description)
	assert.Equal(t, hotelDto.Country, response.Country)
	assert.Equal(t, hotelDto.City, response.City)
	assert.Equal(t, hotelDto.Adress, response.Adress)
	assert.Equal(t, hotelDto.Images, response.Images)
	assert.Equal(t, hotelDto.Amenities, response.Amenities)

	// Asegurarse de que todas las expectativas del mock se cumplan
	mockService.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockService := new(MockHotelService)
	service.HotelService = mockService

	// Datos del hotel actualizado
	updatedHotel := dto.HotelDto{
		Id:          "example-id",
		Name:        "Updated Hotel",
		Description: "Updated Description",
		Country:     "Updated Country",
		City:        "Updated City",
		Adress:      "Updated Address",
		Images:      []string{"updated_image1.jpg", "updated_image2.jpg"},
		Amenities:   []string{"wifi", "pool"},
	}

	// Mockear el servicio para devolver un error
	mockService.On("UpdateHotel", "example-id", updatedHotel).Return(dto.HotelDto{}, errors.NewBadRequestApiError("some error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "example-id"}}
	c.Request, _ = http.NewRequest("PUT", "/hotels/example-id", strings.NewReader(`{
		"Id": "example-id",
		"Name": "Updated Hotel",
		"Description": "Updated Description",
		"Country": "Updated Country",
		"City": "Updated City",
		"Address": "Updated Address",
		"Images": ["updated_image1.jpg", "updated_image2.jpg"],
		"Amenities": ["wifi", "pool"]
	}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.Update(c)

	// Verificar el código de estado y el manejo del error
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar el cuerpo de la respuesta para asegurar que sea un error JSON válido
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON response: %v", err)
	}
	//assert.Contains(t, response["error"], "some error") // Ajustar este mensaje de error según lo que devuelva tu servicio

	mockService.AssertExpectations(t)
}
