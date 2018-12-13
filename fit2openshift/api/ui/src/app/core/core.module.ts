import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';
import { ClarityModule, ClrFormsNextModule } from '@clr/angular';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    HttpClientModule,
    ClarityModule,
    ClrFormsNextModule,
    BrowserAnimationsModule,
  ],
  exports: [
    BrowserModule,
    FormsModule,
    HttpClientModule,
    RouterModule,
    ClarityModule,
    ClrFormsNextModule,
    BrowserAnimationsModule
  ]
})
export class CoreModule { }
