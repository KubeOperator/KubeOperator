import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RegistryComponent } from './registry.component';

describe('RegistryComponent', () => {
  let component: RegistryComponent;
  let fixture: ComponentFixture<RegistryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RegistryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RegistryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
