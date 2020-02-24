import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClarityModule} from '@clr/angular';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientJsonpModule, HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BrowserModule} from '@angular/platform-browser';
import {RouterModule} from '@angular/router';
import {NgxEchartsModule} from 'ngx-echarts';
import {LSelect2Module} from 'ngx-select2';


@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    ClarityModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    HttpClientJsonpModule,
    BrowserAnimationsModule,
    BrowserModule,
    RouterModule,
    NgxEchartsModule,
    LSelect2Module,
  ],
  exports: [
    CommonModule,
    ClarityModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    HttpClientModule,
    HttpClientJsonpModule,
    BrowserModule,
    NgxEchartsModule,
    BrowserAnimationsModule,
    LSelect2Module
  ]
})
export class CoreModule {
}
