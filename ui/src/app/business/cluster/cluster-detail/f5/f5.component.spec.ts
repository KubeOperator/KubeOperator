import { ComponentFixture, TestBed } from '@angular/core/testing';

import { F5Component } from './f5.component';

describe('F5Component', () => {
  let component: F5Component;
  let fixture: ComponentFixture<F5Component>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ F5Component ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(F5Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
