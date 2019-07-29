import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppComponent } from './app.component';
import {NgMatGrpcFormModule} from "ng-mat-grpc-form";

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    NgMatGrpcFormModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
