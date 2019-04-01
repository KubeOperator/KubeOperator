import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClarityModule, ClrFormsDeprecatedModule} from '@clr/angular';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BrowserModule} from '@angular/platform-browser';
import {RouterModule} from '@angular/router';
import {NgxEchartsModule} from 'ngx-echarts';
import {BaseModule} from '../base/base.module';
import {TipModule} from '../tip/tip.module';


@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    ClarityModule,
    ClrFormsDeprecatedModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    BrowserAnimationsModule,
    BrowserModule,
    RouterModule,
    NgxEchartsModule,
  ],
  exports: [
    CommonModule,
    ClarityModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    HttpClientModule,
    ClrFormsDeprecatedModule,
    BrowserModule,
    NgxEchartsModule,
    BrowserAnimationsModule,
  ]
})
export class CoreModule {
}
