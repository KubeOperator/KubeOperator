import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SecurityTaskDeleteComponent } from './security-task-delete.component';

describe('SecurityTaskDeleteComponent', () => {
  let component: SecurityTaskDeleteComponent;
  let fixture: ComponentFixture<SecurityTaskDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecurityTaskDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecurityTaskDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
