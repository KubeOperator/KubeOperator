import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NamespaceDeleteComponent } from './namespace-delete.component';

describe('NamespaceDeleteComponent', () => {
  let component: NamespaceDeleteComponent;
  let fixture: ComponentFixture<NamespaceDeleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ NamespaceDeleteComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NamespaceDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
