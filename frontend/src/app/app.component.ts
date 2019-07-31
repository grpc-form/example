import { Component } from '@angular/core';

@Component({
	selector: 'app-root',
	template: `
<ng-mat-grpc-form name="car" host="http://localhost:8080" (success)="onSuccess($event)"></ng-mat-grpc-form>
  `,
	styles: [`
.mat-angoform-wrapper {
    min-width: 150px;
    max-width: 500px;
    width: 100%;
    margin-bottom: 20px;
}

::ng-deep .mat-form-field {
    width: 100%;
    margin-bottom: -10px;
}

::ng-deep .mat-slider-thumb-label {
    transform: rotate(45deg) !important;
    border-radius: 50% 50% 0 !important;
    background-color: transparent !important;
}

::ng-deep .mat-slider {
    width: 100%;
    margin-top: 10px;
}
::ng-deep .mat-slider-thumb {
    transform: scale(0.5) !important;
}

::ng-deep .mat-slider-thumb-label-text {
    opacity: 1 !important;
}

::ng-deep .mat-radio-container {
    margin-top: 10px !important;
    margin-bottom: 10px !important;
}

::ng-deep .mat-radio-button {
    margin-right: 10px !important;
}
  `]
})
export class AppComponent {
	onSuccess(msg: string) {
		alert(msg);
	}
}
