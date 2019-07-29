package main

import (
	"context"

	"google.golang.org/grpc"

	"example/backend/api"

	grpcform "github.com/grpc-form/api/go"
)

const (
	USERNAME = iota
	AGE
	CAR
	BRAND
)

var (
	NO         int64 = 1
	YES        int64 = 2
	VOLKSWAGEN int64 = 27
)

func main() {
	conn, err := grpc.Dial("database-service:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	db := api.NewDatabaseClient(conn)

	s := grpcform.New()
	s.Add(func() *grpcform.Form {
		return &grpcform.Form{
			Name: "car",
			Fields: []*grpcform.Field{
				USERNAME: {
					InstantValidate: true,
					Label:           "Please enter your username",
					Status:          grpcform.FieldStatus_FIELD_STATUS_REQUIRED,
					TextField: &grpcform.TextField{
						Min:        5,
						MinError:   "Too Short",
						Max:        25,
						MaxError:   "Too long",
						Regex:      "^[a-zA-Z]{5,25}$",
						RegexError: "Only letters",
					},
				},
				AGE: {
					Label:           "How old are you?",
					InstantValidate: true,
					Status:          grpcform.FieldStatus_FIELD_STATUS_REQUIRED,
					NumericField: &grpcform.NumericField{
						Min:   14,
						Max:   99,
						Step:  1,
						Value: 25,
					},
				},
				CAR: {
					Label:           "Do you have a car?",
					InstantValidate: true,
					Status:          grpcform.FieldStatus_FIELD_STATUS_REQUIRED,
					SelectField: &grpcform.SelectField{
						Index: NO,
						Type:  grpcform.SelectType_SELECT_TYPE_SIMPLE,
						Options: []*grpcform.Option{
							{Index: NO, Value: "No"},
							{Index: YES, Value: "Yes"},
						},
					},
					HiddenIf: &grpcform.HiddenIf{
						Validators: []*grpcform.Validator{
							{
								Index:             AGE,
								NumberSmallerThan: 18,
							},
						},
					},
				},
				BRAND: {
					Status:          grpcform.FieldStatus_FIELD_STATUS_HIDDEN,
					InstantValidate: true,
					Label:           "What kind of car are you driving?",
					SelectField: &grpcform.SelectField{
						Index: VOLKSWAGEN,
						Type:  grpcform.SelectType_SELECT_TYPE_MULTI,
						Options: []*grpcform.Option{
							{Index: 1, Value: "Audi"},
							{Index: 2, Value: "BMW"},
							{Index: 3, Value: "Bentley"},
							{Index: 4, Value: "Chevrolet"},
							{Index: 5, Value: "FIAT"},
							{Index: 6, Value: "Ferrari"},
							{Index: 7, Value: "Ford"},
							{Index: 8, Value: "HUMMER"},
							{Index: 9, Value: "Hyundai"},
							{Index: 10, Value: "Jaguar"},
							{Index: 11, Value: "Jeep"},
							{Index: 12, Value: "Kia"},
							{Index: 13, Value: "Lamborghini"},
							{Index: 14, Value: "Maybach"},
							{Index: 15, Value: "Mazda"},
							{Index: 16, Value: "McLaren"},
							{Index: 17, Value: "Mercedes-Benz"},
							{Index: 18, Value: "Nissan"},
							{Index: 19, Value: "Opel"},
							{Index: 20, Value: "Jaguar"},
							{Index: 21, Value: "Porsche"},
							{Index: 22, Value: "Renault"},
							{Index: 23, Value: "Skoda"},
							{Index: 24, Value: "Suzuki"},
							{Index: 25, Value: "Tesla"},
							{Index: 26, Value: "Toyota"},
							{Index: VOLKSWAGEN, Value: "Volkswagen"},
							{Index: 28, Value: "Volvo"},
							{Index: 29, Value: "smart"},
						},
					},
					RequiredIf: &grpcform.RequiredIf{
						Validators: []*grpcform.Validator{
							{
								Index:             AGE,
								NumberGreaterThan: 17,
							},
							{
								Index:         CAR,
								NumberIsEqual: YES,
							},
						},
					},
					HiddenIf: &grpcform.HiddenIf{
						Validators: []*grpcform.Validator{
							{
								Index:             AGE,
								NumberSmallerThan: 18,
							},
							{
								Index:         CAR,
								NumberIsEqual: NO,
							},
						},
					},
				},
			},
			Buttons: []*grpcform.Button{
				{
					Label:  "Reset",
					Type:   grpcform.ButtonFuncType_BUTTON_FUNC_RESET,
					Status: grpcform.ButtonStatus_BUTTON_ACTIVE,
				},
				{
					Label:  "Send",
					Type:   grpcform.ButtonFuncType_BUTTON_FUNC_SEND,
					Status: grpcform.ButtonStatus_BUTTON_DISABLED,
				},
			},
			Valid: false,
		}
	}, func(ctx context.Context, in *grpcform.Form) (res *grpcform.SendFormResponse, err error) {
		out, err := s.ValidateForm(ctx, in)
		if out == nil || err != nil || !out.GetValid() {
			return &grpcform.SendFormResponse{
				Form:    out,
				Succeed: false,
				Message: "Insert Failed: Not Valid",
			}, nil
		}
		index := out.GetFields()[BRAND].GetSelectField().GetIndex()
		options := out.GetFields()[BRAND].GetSelectField().GetOptions()
		var car string
		for _, o := range options {
			if o.GetIndex() == index {
				car = o.GetValue()
			}
		}
		u, err := db.InsertUser(ctx, &api.User{
			Name: out.GetFields()[USERNAME].GetTextField().GetValue(),
			Age:  out.GetFields()[AGE].GetNumericField().GetValue(),
			Car:  car,
		})
		if u == nil {
			return &grpcform.SendFormResponse{
				Form:    out,
				Succeed: false,
				Message: "Insert Failed: DB Error",
			}, nil
		}
		if u.GetName() != out.GetFields()[USERNAME].GetTextField().GetValue() {
			out.GetFields()[USERNAME].Error = "You have already completed this form"
			return &grpcform.SendFormResponse{
				Form:    out,
				Succeed: false,
				Message: "Insert Failed: You have already completed this form",
			}, nil
		}
		return &grpcform.SendFormResponse{
			Form:    out,
			Succeed: true,
			Message: "Thank you!",
		}, nil
	})
	err = s.Start(":50051")
	conn.Close()
	if err != nil {
		panic(err)
	}
}
